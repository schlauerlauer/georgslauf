<?php
require('../session/session.php');
include_once '../host/settings.php';
if(isset($login_session) && $_SESSION['rolle'] == 1) {
  if ($stmt = $mysqli->prepare("SELECT id, kurz, name, stamm, x_axis, y_axis, kategorie, beschreibung, kontakt, anzahl, veggie, material, ort, sonstiges, stufe, color, startGruppen FROM posten WHERE kurz = ? LIMIT 1")) {
    $stmt->bind_param('s', $login_session);
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($id, $kurz, $name, $stamm, $x_axis, $y_axis, $kat, $desc, $kontakt, $anzahl, $veggie, $material, $ort, $sonstiges, $stufe, $color, $startGruppen);
    while ($stmt->fetch()) {
      echo '<div class="ui-bar ui-bar-a"><h3>'.$kurz.' - '.$name.' ('.$stamm.')</h3></div>
      <div class="ui-body ui-body-a">';
      echo '<br><h3 class="ui-bar ui-bar-a ui-corner-all" style="background-color:'.$Farben[$color].';" align="center"><span style="background-color:white;">Postenfarbe</span></h3><br>Beschreibung <strong>'.$desc.'</strong><br><br>Kontakt <strong>'.$kontakt.'</strong><br>Leiteranzahl <strong>'.$anzahl.'</strong> (davon <strong>'.$veggie.'</strong> Veggies)<br><br />';
      if($material != null) echo 'Material '.$material.'<br />';
      if($sonstiges != null) echo 'Sonstiges '.$sonstiges;
      echo "<br />Start Gruppen <strong>".$startGruppen.'</strong>';
      echo '</div>';
    }
  }
}
?>
</div>

<h4 class="ui-bar ui-bar-a ui-corner-all">
    <u><strong>Unser Zeitplan:</strong></u><br>
    <ul>
        <li>8:00 Uhr: Treffen bei Sankt Ansgar </li>
        <li>8:30 Uhr: Auftakt & Andacht</li>
        <li>10:00 Uhr: Startschuss</li>
        <li>13:00 bis 13:30 Uhr: Mittagspause</li>
        <li>16:00 Uhr: Ende des Postenlaufs</li>
        <li>18:00 Uhr: Sieger*innen Ehrung</li>
        <li>circa 20:30 Uhr: Ende</li>
        <li>circa 22:00 Uhr: Party für Leiter*innen und Rover*innen ab 16 Jahren</li>
    </ul>
    <br>
</h4>
<h3 class="ui-bar ui-bar-a ui-corner-all" style="background-color:Tomato; color:white;" align="center">Fehler / Fragen / Hilfe<a data-icon="slack" href="https://bezirkmueisar.slack.com/archives/C03KCKD0DKR" target="_blank">Slack Channel im Bezirksslack</a><a data-icon="phone" href="tel:015756456883">Anrufen (015787297182)</a><a target="_blank" data-icon="comment" href="https://api.whatsapp.com/send?phone=4915787297182&text=Frage%20von%20Posten%20<?php echo $login_session; ?>%3A%0D%0A%0D%0A">WhatsApp (015787297182)</a></h3>
<br>
<br>
<p align="center"><img src="/res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
