<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/includes/connect_gl.php';
include_once 'host/settings.php';
?>
<br>
<?php if($Posten[0] == false) : ?>
<h3 class="ui-bar ui-bar-a ui-corner-all" align="center"><?php echo $Posten[1]; ?></h3>
<br>
<div class="ui-body ui-body-a ui-corner-all" align="center">
  <p>Hier findet Ihr am Tag des Georgslaufs genauere Angaben über die Posten sowie deren exakten Standorte.</p>
</div>
<?php else : ?>
  <div data-role="collapsible-set" data-theme="a">
  	<?php
  		if ($stmt = $mysqli->prepare("SELECT kurz, name, stamm, stufe, x_axis, y_axis FROM posten ORDER BY kurz ASC")) {
  			$stmt->execute();
  			$stmt->store_result();
  			$stmt->bind_result($kurz, $name, $stamm, $stufe, $x_axis, $y_axis);
  			while ($stmt->fetch()) {
          echo '<div data-role="collapsible"><h3>'.$kurz.' - '.$stufe.' - '.$stamm.'</h3><p>';
					echo '<input class="map" p="'.$kurz.'" data-icon="location" x="'.$x_axis.'" y="'.$y_axis.'" type="button" value="Karte"/>';
					echo '</p></div>';
  			}
  		}
  	?>
  </div>
<?php endif; ?>
<br>
<br>
<p align="center"><img src="../res/Schwabing.jpg" style="max-width: 100%; max-height: 700px;"/></p>
<div align="center">Georgslauf 2015 - im schönen Schwabing</div>
