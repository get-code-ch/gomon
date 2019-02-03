//printjson(db.adminCommand('listDatabases'));
db = db.getSiblingDB('gomondb');
//
// menu collection
//
db.menu.drop();
db.createCollection("menu");
db.menu.insert({"Key": "100.HOME", "Link": "/", "Text": "Home", "Visible": true});
db.menu.insert({"Key": "201.HOSTS", "Link": "/hosts", "Text": "Hosts", "Visible": true});
db.menu.insert({"Key": "202.PROBE", "Link": "/probes", "Text": "Probe", "Visible": true});
db.menu.insert({"Key": "900.LOGOUT", "Link": "/logout", "Text": "Logout", "Visible": true});
db.menu.insert({"Key": "800.ADMIN", "Link": "/admin", "Text": "Admin", "Visible": true});
//
// status collection
//
db.status.drop();
db.createCollection("status");
db.status.insert({"Key": "N/A", "Description": "N/A"});
db.status.insert({"Key": "Ok", "Description": "Ok", "Color": "#007700"});
db.status.insert({"Key": "Warning", "Description": "Warning", "Color": "#770000"});
db.status.insert({"Key": "Error", "Description": "Warning", "Color": "#aa4400"});
