<?php
//Einstellungsdatei für den Georgslauf

$Email = "gl22@stamm-prm.de";

$Kat = array('Wissen & Technik','Spiel & Spaß'); //Die Kategorien für die Posten, beliebig erweiterbar / veränderbar!
$Stufe = array('Wölflinge','Jupfis','Pfadis','Rover'); // Stufen namen, wie sie angezeigt werden (auf anmeldeseite aber hardcoded!) (in SQL sind nur nummern) eher nicht verändern

$WKat = array(9,10); //WöPo-Anzahl Posten pro Kategorie - Anzahl an Zahlen muss mit Kategorien ($Kat) Anzahl übereinstimmen!
$RKat = array(2,2,); //RoPo-Anzahl Posten pro Kategorie - "


$Anmeldung = true; //Stämme können Gruppen und Posten anmelden
$Abmeldung = false; //Stämme können Gruppen und Posten löschen
$PAnmeldung = true; //Stämme können Posten anmelden
$GAnmeldung = true; //Stämme können Gruppen anmelden
$ShowLogin = false; //Show Posten Login Information in Stamm -> Posten

$Rollen = array("keine Rechte", "Posten", "Stamm", "Host", "Admin"); //Rechte vergabe

$Posten = array(false, "Die Posten sind noch nicht öffentlich"); // Posten auf Hauptseite sichtbar oder nicht
$PostenStufe = 'Z'; //Trennung RoPo / WöPo großes o, -> siehe Posten -> update.php
$Punkte = array(true, "Die Posten dürfen noch keine Punkte vergeben!"); //true = Posten dürfen Punkte verteilen, false -> Fehlermeldung;

$Ende = false; //Host kann keine Punkte mehr eintragen wenn true;
$Backup = true; //Erstelle automatische Backups der Datenbank;

$Farben = array('LimeGreen','MediumVioletRed','yellow','SkyBlue'); //Postenfarben in HTML Color Names
?>
