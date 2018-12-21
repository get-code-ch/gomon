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
//
// Probes collection
//
db.probes.drop();
db.createCollection("probes");
db.probes.insert({
    "Key": "WPPAGE",
    "Name": "WordPress pages probe",
    "Description": "Probe if WordPress site is up and running"
});
db.probes.insert({
    "Key": "PING",
    "Name": "Basic ping probe",
    "Description": "If IP address respond of ping",
    "Command": "ping {ip}"
});
//
// Hosts collection
//
db.hosts.drop();
db.createCollection("hosts");
db.hosts.insert({"Key": "GETCODEWWW", "Name": "www.get-code.ch", "FQDN": "www.get-code.ch"});
db.hosts.insert({"Key": "GOOGLEDNS", "Name": "Google DNS", "FQDN": "google-public-dns-a.google.com", "IP": "8.8.8.8"});
//
// States collection
//
db.states.drop();
db.createCollection("states");
db.states.insert({"Host": "GETCODEWWW", "Probe": "PING", "Status": "N/A"});
print("Collections created!");
