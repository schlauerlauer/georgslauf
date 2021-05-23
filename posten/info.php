<?php
require('../session/session.php');
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../host/settings.php';
?>
<h3 class="ui-bar ui-bar-a ui-corner-all" align="center">Posteninfo</h3>
<br>
<div class="ui-corner-all">
<?php
if(isset($login_session) && $_SESSION['rolle'] == 1) {
  if ($stmt = $mysqli->prepare("SELECT id, kurz, name, stamm, x_axis, y_axis, kategorie, beschreibung, kontakt, anzahl, veggie, material, ort, sonstiges, stufe, color, startGruppen FROM posten WHERE kurz = ? LIMIT 1")) {
    $stmt->bind_param('s', $login_session);
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($id, $kurz, $name, $stamm, $x_axis, $y_axis, $kat, $desc, $kontakt, $anzahl, $veggie, $material, $ort, $sonstiges, $stufe, $color, $startGruppen);
    while ($stmt->fetch()) {
      echo '<div class="ui-bar ui-bar-a"><h3>'.$kurz.' - '.$name.' ('.$stamm.')</h3></div>
      <div class="ui-body ui-body-a">Posten für <strong>';
      if ($stufe == 'RoPo') echo 'Pfadis & Rover  (RoPo)';
      else echo 'Wölflinge & Jupfis  (WöPo)';
      echo '</strong><br><h3 class="ui-bar ui-bar-a ui-corner-all" style="background-color:'.$Farben[$color].';" align="center"><span style="background-color:white;">Postenfarbe</span></h3><br>Beschreibung <strong>'.$desc.'</strong><br><br>Kontakt <strong>'.$kontakt.'</strong><br>Leiteranzahl <strong>'.$anzahl.'</strong> (davon <strong>'.$veggie.'</strong> Veggies)<br><br />';
      if($material != null) echo 'Material '.$material.'<br />';
      if($sonstiges != null) echo 'Sonstiges '.$sonstiges;
      echo "<br />Start Gruppen <strong>".$startGruppen.'</strong>';
      echo '</div>';
    }
  }
}
?>
</div>
<br>
<input class="map" p="<?php echo $kurz; ?>" data-icon="location" x="<?php echo $x_axis; ?>" y="<?php echo $y_axis; ?>" type="button" value="Position"/>
<br />
<h4 class="ui-bar ui-bar-a ui-corner-all">
<br>
Unser Zeitplan für den 11.07.2021 sieht dann folgendermaßen aus:<br>
<ul>
  <li>08:00 Uhr: Treffen in euren Stämmen</li>
  <li>09:00 Uhr: Spiritueller Auftakt</li>
  <li>10:30 Uhr: Startschuss</li>
  <li>12:30 bis 13:00 Uhr: Mittagspause</li>
  <li>16:30 Uhr: Ende des Laufs</li>
  <li>18:00 Uhr: Siegerehrung</li>
  <li>circa 20:30 Uhr: Ende</li> 
  <li>circa 21:30 Uhr: Party</li> 
</ul>
<br>
</h4>
<h3 class="ui-bar ui-bar-a ui-corner-all" style="background-color:Tomato; color:white;" align="center">Fehler / Fragen / Hilfe<a data-icon="mail" href="mailto:gl21@stamm-prm.de">Mail an gl21@stamm-prm.de</a><a data-icon="phone" href="tel:015756456883">Anrufen (015756456883)</a><a target="_blank" data-icon="comment" href="https://api.whatsapp.com/send?phone=4915756456883&text=Frage%20von%20Posten%20<?php echo $login_session; ?>%3A%0D%0A%0D%0A">WhatsApp (015756456883)</a></h3>
<br>
<br>
<p align="center"><img src="/res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
