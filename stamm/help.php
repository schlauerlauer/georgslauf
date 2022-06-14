<?php
if (isset($_POST['help'])) {
	switch($_POST['help']) {
		//Die Überschrift des Hilfs-Pop-Ups lautet "Hilfe zu den ... " (Deshalb "GruppennameN")
		case "Gruppenstufen":
			echo "Mischgruppen als Stufe der Ältesten Kinder anmelden.";
				break;
		case "Gruppennamen":
			echo "Euer Gruppenname sollte einzigartig sein.<br>
			<br>
			Der Gruppenname wird benötigt!
			";
			break;
		case "Gruppengrößen":
			echo "Generell sind Gruppengrößen nur von 3 bis 12 Kindern erlaubt. <br> Teilen der Gruppen erst bei 13 Kindern erlaubt in zwei Laufgruppen.
			<br>
			Für Ausnahmen bitte uns kontaktieren, wir werden <strong>versuchen</strong> diese zu ermöglichen.";
			break;
		case "Kategorien":
			echo "Eure Posten sollen zu den unten stehenden Kategorien passen.<br>
			<br>
			Die mögliche Postenanzahl: Wissen & Technik (9x), Spiel & Spaß (10x)";
			break;
		case "Stufen":
			echo "Eure Posten werden besucht von entweder<br>
			Wölflingen & Jupfis <strong>oder</strong> Pfadis & Rovern!<br>
			<br>
			Passt euer Programm dementsprechend an.<br>
			Jeder Stamm kann nur einen RoPo stellen!<br>
			Die anderen Posten müssen WöPos sein.<br>
			<br>
			Wenn ihr nur einen Posten stellen könnt dann eher einen WöPo.";
			break;
		case "Nummern":
			echo "Bei jüngeren Laufgruppen mit Leiter,<br>
			bitte die Handynummer des Leiters eintragen<br>
			(Für die Whatsappgruppe).<br>
			Hilfreich wäre mit Vorwahl (+49151...).";
			break;
		default:
			echo "Keine ID angegeben";
			break;
	}
} else echo "Es ist ein Fehler aufgetreten.";
?>
