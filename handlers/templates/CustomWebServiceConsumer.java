package com.datacert.apps.mattermanagement.integration;

import groovy.lang.GroovyInterceptable;
import groovy.xml.*;
import java.util.List;
import java.util.Map;
import java.io.File;
import java.text.SimpleDateFormat;
import com.datacert.core.util.Logger;
import groovy.sql.Sql;

import org.springframework.stereotype.Component;


/**
 * Data Provider for Beta web service integration. 
 */
@Component
public class CustomWebServiceConsumer extends IDataProvider 
	implements GroovyInterceptable{

	/** logger */
	private def logger;
	/** a map for configurations from the UI */
	private def configurationMap
	/** batch size */
	private def batchSize;

	private def jarDirectory;
	private def webServiceEndPoint;
	private def webServiceSoapAction
	private def maxNumRecord;
	private searchBy;
	private searchValue;
	private timeStamp;


	/**
	 * default constructor
	 */
	public CustomBetaSoapConsumer(){}

	
	/**
	 * A constructor for initializing the logger, batch size and configuration map
	 *
	 * @param logger				logger
	 * @param batchSize				batch size
	 * @param configurationMap		configuration map
	 */
	public CustomBetaSoapConsumer(Logger logger, long batchSize, Map configurationMap) {
		this.logger = logger;
		this.batchSize = batchSize;
		this.configurationMap = configurationMap;
	}
	
	
	/**
	 * Check configurationMap input.  
	 *
	 */
	public boolean initialize() throws Exception {

		jarDirectory = configurationMap.get("JarDirectory");
		logger.debug('jarDirectory value is ' + jarDirectory);
		webServiceEndPoint = configurationMap.get("WebServiceEndPoint");
		logger.debug('webServiceEndPoint value is ' + webServiceEndPoint);
		webServiceSoapAction = configurationMap.get("WebServiceSoapAction");
		logger.debug('webServiceSoapAction value is ' + webServiceSoapAction);
		maxNumRecord = configurationMap.get("MaxNumRecord");
		logger.debug('maxNumRecord value is ' + maxNumRecord);
		def sdf = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss");
		timeStamp = sdf.format(new Date());
		logger.debug('timeStamp value is ' + timeStamp);
		searchBy = configurationMap.get("SearchBy");
		logger.debug('searchBy value is ' + searchBy);
		searchValue = configurationMap.get("SearchValue");
		logger.debug('searchValue value is ' + searchValue);
		
		
		return true;
	}

	
	/**
	 *
	 * @return		All records to be created/updated.
	 * @throws		Exception
	 */
	public Map<String, List<?>> getNextBatch() throws Exception {
	
		Map<String, List<?>> dataMap = new HashMap<String, List<?>>() 

		/**
		client is hosting Passport in their own environment
		instead of placing the WSLite jar in the lib directory under tomcat,
		we will be placing in a custom location so that it persists during upgrades
		as a result, we are loading the classes manually, as opposed to using import statements
		*/
		def loader = this.class.classLoader
		def url = new File(jarDirectory + "groovy-wslite-0.7.1.jar").toURI().toURL();
		logger.debug('url is : ' + url);
		loader.addURL(url)
		def soapClass = Class.forName("wslite.soap.SOAPClient")
		
		/**
		client does not require authentication, so we can skip loading HTTPBasicAuthorization
		//def httpClass = Class.forName("wslite.http.auth.HTTPBasicAuthorization")
		*/
		def client = soapClass.newInstance(webServiceEndPoint);

		/**
		we are passing these parameters to a utility method that in turn passes them as part of the soap request
		response is returned
		*/
		def payload = getXml(searchBy,searchValue,maxNumRecord,timeStamp);
		def response = client.send(SOAPAction:webServiceSoapAction,payload)
		
		/**
			success value of web service call will be 'success','failure', or 'No Record Found'
		*/
		def successResult = response.getAccountInfoResponse.responseMessageHeader.success
		
		if(successResult == 'success'){
			/**
				customer details 
			*/
			def customerDetails = response.getAccountInfoResponse.AccountInfo.CustomerDetails;
			/**
				account details associated with customer
			*/
			def accountDetails = response.getAccountInfoResponse.AccountInfo.account;
		  dataMap.put('1',[
			  'person':[
			  customerDetails.FirstName,
			  customerDetails.LastName,
			  customerDetails.MiddleName,
			  customerDetails.Salutation,
			  customerDetails.Suffix,
			  customerDetails.SSNTIN,
			  customerDetails.HoldersBirthDate,
			  customerDetails.HomePhone,
			  customerDetails.BusinessPhone,
			  customerDetails.AddressType,
			  customerDetails.AddressLine1,
			  customerDetails.AddressLine2,
			  customerDetails.City,
			  customerDetails.State,
			  customerDetails.Zip
			  ],'accounts':getAccountsList(accountDetails)
		  ]);
		}else if(successResult == 'failure'){
			def errorDetails = response.getAccountInfoResponse.responseMessageHeader.fault;
			dataMap.put('result','failure');
			dataMap.put('errorCode', errorDetails.errorCode);
			dataMap.put('errorMessage',errorDetails.errorMessage);
		}else if(successResult == 'No Record found'){
			dataMap.put('result','No Record found');
		}else{
			dataMap.put('result',successResult);
		}

		return dataMap;
	}
	
	def getAccountsList(def personAccounts){
		def resultList = [];
		def resultMap;
		personAccounts.each{
			resultMap = new HashMap();
			def currentAccount = []
			currentAccount.add(it.SSNTIN)
			currentAccount.add(it.AccountNo)
			currentAccount.add(it.AccountDescription)
			currentAccount.add(it.InvestmentObjective)
			currentAccount.add(it.AccountStatus)
			currentAccount.add(it.LastChangeDate)
			currentAccount.add(it.AcctClass)
			currentAccount.add(it.AccountName)
			resultList.add(currentAccount);
		}
		return resultList;
	}

	
	def getXml(searchBy,searchValue,maxNumRecord,timeStamp){
		return """
			<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
			xmlns:ns="http://lpl.com/BTS/BetaAccount/getBetaAccountInfoServiceTypes/2014/08" 
			xmlns:ns1="http://lpl.com/BTS/Framework/MessageHeader/2013/08" 
			xmlns:ns2="http://lpl.com/BTS/BetaAccount/getBetaAccountInfoRequest/2014/08">
			   <soapenv:Header/>
			   <soapenv:Body>
				  <ns:getAccountInfoRequest>
					 <ns1:requestMessageHeader>
						<ns1:securityToken></ns1:securityToken>
						<ns1:organization></ns1:organization>
						<ns1:transactionId></ns1:transactionId>
						<ns1:timeStamp>${timeStamp}</ns1:timeStamp>
						<ns1:sourceSystem></ns1:sourceSystem>
						<ns1:hostName></ns1:hostName>
						<!--Optional:-->
						<ns1:userId></ns1:userId>
					 </ns1:requestMessageHeader>
					 <ns2:AccountInfoReq>
						<ns2:SearchBy>${searchBy}</ns2:SearchBy>
						<ns2:SearchValue>${searchValue}</ns2:SearchValue>
						<ns2:MaxNumRecord>${maxNumRecord}</ns2:MaxNumRecord>
					 </ns2:AccountInfoReq>
				  </ns:getAccountInfoRequest>
			   </soapenv:Body>
			</soapenv:Envelope>
		"""
	}	

	/**
	 *
	 */
	public def complete(){
		logger.debug("Inside complete()");
	}
	
 
}