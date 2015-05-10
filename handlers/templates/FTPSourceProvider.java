package com.datacert.apps.mattermanagement.integration;

import org.apache.commons.vfs2.auth.StaticUserAuthenticator;
import org.apache.commons.vfs2.FileSystemOptions;
import org.apache.commons.vfs2.impl.DefaultFileSystemConfigBuilder;
import org.apache.commons.vfs2.FileType;
import org.apache.commons.vfs2.FileSystemException;
import org.apache.commons.vfs2.AllFileSelector;
import org.apache.commons.vfs2.Selectors;
import com.datacert.core.util.Logger;

/*
Configurations 
--------------
ServerName	: FTP / SFTP Server name or IP address
UserName	: User Name
Password	: Password
FileSystem	: protocol to be used (ftp|sftp)
PortNumber	: server port number
Source		: 'Local' for local to remote ; 'Remote' for remote to> local
RemoteDir	: FTP directory from / to which files have to be moved
LocalDir	: local directory from / to which files have to be moved
FileMask	: file mask
ArchiveOnFTP	: 'true' if you want to archive files in FTP
ArchiveOnLocal	: 'true' if you want to archive files in local system
ArchiveDir	: acrhive dir if ArchiveOnFTP=true or ArchiveOnLocal=true

*/

public class FTPSourceProvider extends IExternalProvider {

	private def logger;
	private def configurationMap
	private def fileList;
	private def fsManager;
	private def options;
	private def sftpFile;
	private String startPath;

	public FTPSourceProvider(){
	}

  public FTPSourceProvider(Logger logger){
	}
	public FTPSourceProvider(Logger logger, Map configurationMap) {
		this.logger = logger;
		this.configurationMap = configurationMap;
	}

	public boolean initialize() {
		try {
			if(!configurationMap.get("ServerName")) throw new RuntimeException("Please set value for ServerName");
			if(!configurationMap.get("UserName")) throw new RuntimeException( "Please set value for UserName");
			if(!configurationMap.get("Password")) throw new RuntimeException( "Please set value for Password");
			if(!configurationMap.get("RemoteDir")) throw new RuntimeException( "Please set value for RemoteDir");
			if(!configurationMap.get("FileMask")) throw new RuntimeException( "Please set value for FileMask");
			if(!configurationMap.get("FileSystem")) throw new RuntimeException( "Please set value for FileSystem");
			if(!configurationMap.get("PortNumber")) throw new RuntimeException( "Please set value for PortNumber");
			if(!configurationMap.get("LocalDir")) throw new RuntimeException( "Please set value for LocalDir");
			if(!configurationMap.get("Source")) throw new RuntimeException( "Please set value for Source");
			
			fsManager = org.apache.commons.vfs2.VFS.getManager();

			def auth = new StaticUserAuthenticator(null, configurationMap.get("UserName").toString(),	configurationMap.get("Password").toString());
			options = new FileSystemOptions();
			startPath = configurationMap.get("FileSystem").toString()+ "://" + configurationMap.get("ServerName").toString() + ':' + configurationMap.get("PortNumber").toString();
			DefaultFileSystemConfigBuilder.getInstance().setUserAuthenticator(options,	auth);

			//sftpFile = fsManager.resolveFile(startPath+configurationMap.get("RemoteDir").toString(), options);
			sftpFile = resolveSftpFile(startPath+configurationMap.get("RemoteDir").toString(), options, configurationMap.get("Retry"));
			logger.info(configurationMap.get("FileSystem").toString().toUpperCase() + " connection successfully established to " +startPath);

		}
		catch (Exception ex) {
			logger.error("Error in FTP initialize",ex)
			return false
		}
		return true;
	} 

	private def resolveSftpFile(filePath, options, retry) throws Exception {
		def retryConfig = 0
		if (retry == null) {
			retryConfig = 3	//default
		} else {
			retryConfig = retry.toString().toInteger()
		}
		def retryCount = 0
		while (true) {
			try {
				return fsManager.resolveFile(filePath, options);
			} catch (FileSystemException ex) {
				if (ex.message.contains("Could not connect")) {
					if (retryCount != retryConfig) {
						retryCount++
						logger.info("Retry connecting to file server: ${retryCount} - ${retryConfig} ...")
					} else {
						throw ex
					}
				} else {
					throw ex
				}
			}
		}
		return null
	}
	
	public void process() {
		if(configurationMap.get("Source").toString().equals("Remote")){
			downloadFiles();
			if(configurationMap.get("ArchiveOnFTP")?.toString().equals("true"))
				archiveRemoteFiles();
			if(configurationMap.get("DeleteOnFTP")?.toString().equals("true"))
				deleteRemoteFiles();
		}
		else
		{
			uploadFiles();
			if(configurationMap.get("ArchiveOnLocal")?.toString().equals("true"))
				archiveLocalFiles();
		}
	}

	private void downloadFiles() {
		fileList = new ArrayList<String>();
		def localDir = configurationMap.get("LocalDir").toString()
		sftpFile.children.findAll{it.type == FileType.FILE && it.name.baseName.matches(fileMaskToRegex(configurationMap.get("FileMask")))}.each(){
			f ->
			try{
				def fileName = f.name.baseName;
				def localFile =  fsManager.resolveFile("file://$localDir$File.separator$fileName");
				if (!localFile.getParent().exists()) {
					localFile.getParent().createFolder();
				}
				localFile.copyFrom(f,new AllFileSelector());
				logger.info( fileName + " downloaded to " + localDir);
				fileList.add("$localDir$File.separator$fileName");
			} catch(Exception ex){
				logger.error( "error downloading file $f.name:", ex);
			}
		}
	}
	
	public void uploadFiles() {
		def filePath = new File(configurationMap.get("LocalDir").value?.toString());

		filePath.listFiles().findAll{it.isFile() && it.name.matches(fileMaskToRegex(configurationMap.get("FileMask")))}.each(){
			try{
				def localFile = fsManager.resolveFile("file://$it");
				
				def remoteFile = fsManager.resolveFile(startPath+configurationMap.get("RemoteDir").toString() + File.separator  + it.name, options);
				
				remoteFile.copyFrom(localFile, Selectors.SELECT_SELF);
				logger.info( it.name + " uploaded to " + remoteFile.getName());
			} catch(Exception ex) {
				logger.error( "error uploading file $it:", ex);
			}
		}
	}

	public void archiveRemoteFiles() throws Exception {
		if(!configurationMap.get("ArchiveDir")) throw new RuntimeException( "Please set value for ArchiveDir")
		sftpFile.children.findAll{it.type == FileType.FILE && it.name.baseName.matches(fileMaskToRegex(configurationMap.get("FileMask")))}.each(){
			try{
				def remoteFile = fsManager.resolveFile(startPath+ configurationMap.get('RemoteDir')+configurationMap.get('ArchiveDir') + File.separator+ it.name.baseName, options);
				it. moveTo(remoteFile);
				logger.info( it.name.baseName + " archived");
			} catch (Exception ex) {
				logger.error("Error in FTP Archive",ex)
			}
		}
	}
	
	public void archiveRemoteFailedFiles(def fileNameToArchive) throws Exception {
		if(!configurationMap.get("FailedArchiveDirectory")) throw new RuntimeException( "Please set value for ArchiveDir")
		def sourceFile = sftpFile.resolveFile('Archive/' + fileNameToArchive);
		logger.debug('sourceFile is ' + sourceFile.name.baseName);
		//sftpFile.children.findAll{it.type == FileType.FILE && it.name.baseName.matches(fileMaskToRegex(fileNameToArchive))}.each(){
			try{
				def remoteFile = fsManager.resolveFile(startPath+ configurationMap.get('RemoteDir')+configurationMap.get('FailedArchiveDirectory') + File.separator+ sourceFile.name.baseName, options);
				sourceFile.moveTo(remoteFile);
				logger.info( sourceFile.name.baseName + " moved to failed archive");
			} catch (Exception ex) {
				logger.error("Error in FTP Failed Archive",ex)
			}
		//}
		
	}
	
	public void close(){
		def fs = sftpFile?.getFileSystem();
		this.fsManager?.closeFileSystem(fs);
	}
	
	public void deleteRemoteFiles() throws Exception {
		sftpFile.children.findAll{it.type == FileType.FILE && it.name.baseName.matches(fileMaskToRegex(configurationMap.get("FileMask")))}.each(){
			try{
				it.delete();
				logger.info( it.name.baseName + " deleted");
			} catch (Exception ex) {
				logger.error("Error in deleting file on FTP",ex)
			}
		}
	}
	
	public void archiveLocalFiles() throws Exception {
		if(!configurationMap.get("ArchiveDir")) throw new RuntimeException( "Please set value for ArchiveDir")
		filePath.listFiles().findAll{it.isFile() && it.name.matches(fileMaskToRegex(configurationMap.get("FileMask")))}.each(){
			try{
				def localFile = fsManager.resolveFile("file://$it");
				def arcDir = configurationMap.get('LocalDir')+configurationMap.get('ArchiveDir')
				def archiveFile = fsManager.resolveFile("file://$arcDir$File.separator$it.name");
				localFile.moveTo(archiveFile);
				logger.info( it.name + " archived");
			} catch (Exception ex) {
				logger.error("Error in local Archive",ex)
			}
		}
	}

	public def complete() throws Exception {
/* 		def fs = sftpFile?.getFileSystem();
		this.fsManager?.closeFileSystem(sftpFile?.getFileSystem());
 */		//moved to close() for tt 4852
		return fileList
	}
	
	private def fileMaskToRegex(String fileMask){
		return fileMask.replaceAll("\\.", "[.]").replaceAll("\\*", ".*").replaceAll("\\?", ".");
	}
}