s = SpreadsheetApp.openById("159Y0xb2X3r4QOxv1ccjyF1VBvE5W4pJ5aYV-2trn9jQ");
var sheet = ss.getSheetByName("Data draft All locations");
var data = sheet.getDataRange().getValues();
var headings = data[0];

/* Take a organization name as input and return the
 * row corresponding to that organization name.*/
 
function organizationQuery(organizationName){
 for (var i = 1; i < data.length; i++){
  if (organizationName === data[i][1]){
    return data[i]
  }
 }
}
 
function zipcodeQuery(zipcode) { 
 for (var i = 1; i < data.length; i++){
  if (zipcode === data[i][3].toString()){
    return data[i]
  }
 }
}

function categoryQuery(category){
 for (var i = 1; i < data.length; i++){
  if (category === data[i][0]){
    return data[i]
  }
 }
}


/* Take a spreadsheet (organization) row and turn it into an object
 with spreadsheet headings as object keys. */
 
function formatOrganization(rowData){
 var organization = {}
 for (var i = 0; i < headings.length; i++){
   Logger.log('Headings: ' + headings);
   organization[headings[i].toString()] = rowData[i];
 }
 return  organization
}

/*
function doGet(e) {
  var params = JSON.stringify(e);
  return HtmlService.createHtmlOutput(params);
}
*/

function executeOrganizationNameQuery(request) {
 
     organizationNames = request.parameters.orgName;
 
      // The object to be returned as JSON
      response = {
        organizations : []
      }
  
      // Fill the organzations array with requested organizations
      for (var i = 0; i < organizationNames.length; i++){ 
        sheetData = organizationQuery(organizationNames[i])
        if(sheetData !== undefined) {
          var org = formatOrganization(sheetData)
          if(org !== undefined) {
            response.organizations.push(org)
          }
        }
      }
 
      if (response.organizations.length > 0)
      {
        return ContentService.createTextOutput(JSON.stringify(response));
      } 
      else 
      {
        return ContentService.createTextOutput('Invalid Request. Organization Name(s) not found.');
      } 
  
}

function executeZipcodeQuery(request) {
 
     zipcodes = request.parameters.zipcode;
 
      // The object to be returned as JSON
      response = {
        organizations : []
      }
  
      // Fill the organzations array with requested organizations
      for (var i = 0; i < zipcodes.length; i++){ 
        sheetData = zipcodeQuery(zipcodes[i])
        if(sheetData !== undefined) {
          var org = formatOrganization(sheetData)
          if(org !== undefined) {
            response.organizations.push(org)
          }
        }
      }
 
      if (response.organizations.length > 0)
      {
        return ContentService.createTextOutput(JSON.stringify(response));
      } 
      else 
      {
        return ContentService.createTextOutput('Invalid Request. zipcode(s) not found.');
      } 
  
}

function executeCategoryQuery(request) {
 
     categories = request.parameters.category;
 
      // The object to be returned as JSON
      response = {
        organizations : []
      }
  
      // Fill the organzations array with requested organizations
      for (var i = 0; i < categories.length; i++){ 
        sheetData = categoryQuery(categories[i])
        if(sheetData !== undefined) {
          var org = formatOrganization(sheetData)
          if(org !== undefined) {
            response.organizations.push(org)
          }
        }
      }
 
      if (response.organizations.length > 0)
      {
        return ContentService.createTextOutput(JSON.stringify(response));
      } 
      else 
      {
        return ContentService.createTextOutput('Invalid Request. Category(ies) not found.');
      } 
  
}


function doGet(request) {
  
    // Check for a valid request URI
    if (request.parameter.orgName !== undefined)
    {
        return executeOrganizationNameQuery(request);
    }     
    else if(request.parameter.zipcode !== undefined)
    {
        return executeZipcodeQuery(request); 
    }
    else if(request.parameter.category !== undefined)
    {
        return executeCategoryQuery(request); 
    }
    else 
    {
      return ContentService.createTextOutput('Invalid Request. Use at least one valid "orgname" parameter.');
    }

}  


function testDoGetWithOrgName() {
  var testRequest = {"parameter":{"orgName":"Bethesda Project"},"contextPath":"","contentLength":-1,"queryString":"orgName=Bethesda%20Project","parameters":{"orgName":["Bethesda Project"]}};
  
  doGet(testRequest);
}

function testDoGetWithZipcode() {
  var testRequest = {"parameter":{"zipcode":"19132"},"contextPath":"","contentLength":-1,"queryString":"zipcode=19132","parameters":{"zipcode":["19132"]}};
  
  doGet(testRequest);
}

function testDoGetWithCategory() {
  var testRequest = {"parameter":{"category":"emergency shelter"},"contextPath":"","contentLength":-1,"queryString":"category=emergency%20shelter","parameters":{"category":["emergency shelter"]}};
  
  doGet(testRequest);
}




function logHomelessServicesInfo() {
  var sheet = SpreadsheetApp.getActiveSheet();
  var data = sheet.getDataRange().getValues();
  for (var i = 0; i < data.length; i++) {
    Logger.log('Category: ' + data[i][0]);
    Logger.log('Organization Name: ' + data[i][1]);
    Logger.log('Address: ' + data[i][2]);
  }
}
