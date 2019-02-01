<h1>Einstellungen</h1>
<button class="ui-btn ui-btn-inline ui-icon-back ui-btn-icon-left" onclick="window.location.href='/host'">Zurück</button>
<?php

$data = array();
$data['host'] = "Maxkolbe"; //Stamm der den Georgslauf austrägt! - Nur dieser Name hat Zugriff auf die Host & Einstellungsseite
$data['punkte'] = array(true, "Die Posten dürfen noch keine Punkte vergeben!"); //true = Posten dürfen Punkte verteilen, false -> Fehlermeldung;

$data['kategorie'] = array('Kreativ','Pfadiwissen','Wissen','Action','erste Hilfe'); //Die Kategorien für die Posten, beliebig erweiterbar / veränderbar!
$data['stufen'] = array('Wölflinge','Jupfis','Pfadis','Rover'); // Stufen namen, wie sie angezeigt werden (auf anmeldeseite aber hardcoded!) (in SQL sind nur nummern) eher nicht verändern
$data['wKat'] = array(3,3,2,3,1); //WöPo-Anzahl Posten pro Kategorie - Anzahl an Zahlen muss mit Kategorien ($Kat) Anzahl übereinstimmen!
$data['rKat'] = array(2,2,2,3,1); //RoPo-Anzahl Posten pro Kategorie - "

$data['stämme'] = array("Canisius", "FC", "Hl.Engel", "Hl.Kreuz", "MariaHilf", "Maxkolbe", "PRM", "St.Anna", "St.Ansgar", "St.Severin", "Swapingo");

$file = 'settings.json';
$content = json_encode($data);
file_put_contents($file, $content);
$content = json_decode(file_get_contents($file), TRUE);
?>

<ul data-role="listview" data-inset="true" data-theme="a">
  <label for="host" class="select">Host Stamm</label>
    <select name="host" id="host">
      <?php
      echo '<option value="'.$content['host'].'">'.$content['host'].'</option>';
      for ($i = 0; $i < sizeof($content['stämme']); $i++) {
        if ($content['stämme'][$i] != $content['host']) echo '<option value="'.$content['stämme'][$i].'">'.$content['stämme'][$i].'</option>';
      }
      ?>
    </select>
</ul>
<ul data-role="listview" data-inset="true" data-theme="a">
  <li class="ui-field-contain">
    <label for="flipPunkte">Posten können Gruppen bewerten</label>
    <input data-role="flipswitch" id="flipPunkte" <?php if($content['punkte'][0] == true) echo 'checked=""'; ?> type="checkbox">
  </li>
   <li data-role="rangeslider">
     Gruppengröße
     <input name="range-1a" id="range-1a" min="1" max="15" value="5" type="range">
     <input name="range-1b" id="range-1b" min="1" max="15" value="12" type="range">
   </li>
   <li class="ui-body ui-body-b">
     <div><button class="ui-btn ui-corner-all ui-btn-a">Speichern</button></div>
   </li>
</ul>
