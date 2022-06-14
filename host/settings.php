<?php
//Einstellungsdatei für den Georgslauf

$Email = "gl22@stamm-prm.de";

$Kat = array('Kreativ','Pfadiwissen','Wissen','Action','erste Hilfe'); //Die Kategorien für die Posten, beliebig erweiterbar / veränderbar!
$Stufe = array('Wölflinge','Jupfis','Pfadis','Rover'); // Stufen namen, wie sie angezeigt werden (auf anmeldeseite aber hardcoded!) (in SQL sind nur nummern) eher nicht verändern

$WKat = array(4,2,2,3,1); //WöPo-Anzahl Posten pro Kategorie - Anzahl an Zahlen muss mit Kategorien ($Kat) Anzahl übereinstimmen!
$RKat = array(2,2,2,3,1); //RoPo-Anzahl Posten pro Kategorie - "


$Anmeldung = false; //Stämme können Gruppen und Posten anmelden
$Abmeldung = false; //Stämme können Gruppen und Posten löschen
$PAnmeldung = false; //Stämme können Gruppen anmelden
$GAnmeldung = false; //Stämme können Posten anmelden
$ShowLogin = false; //Show Posten Login Information in Stamm -> Posten

$Rollen = array("keine Rechte", "Posten", "Stamm", "Host", "Admin"); //Rechte vergabe

$Posten = false(true, "Die Posten sind noch nicht öffentlich"); // Posten auf Hauptseite sichtbar oder nicht
$PostenStufe = 'Z'; //Trennung RoPo / WöPo großes o, -> siehe Posten -> update.php
$Punkte = array(true, "Die Posten dürfen noch keine Punkte vergeben!"); //true = Posten dürfen Punkte verteilen, false -> Fehlermeldung;

$Ende = false; //Host kann keine Punkte mehr eintragen wenn true;
$Backup = true; //Erstelle automatische Backups der Datenbank;

$Farben = array('LimeGreen','MediumVioletRed','yellow','SkyBlue'); //Postenfarben in HTML Color Names
?>
