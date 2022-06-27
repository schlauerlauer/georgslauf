<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<h2>Wächterwanze</h2>
<br />
<h3>Bewertung im anderen Lauf</h3>
<ul data-role="listview" data-theme="a" data-inset="true" data-split-icon="delete" data-split-theme="a">
	<?php
			$posten = 0;
			if ($stmt = $mysqli->prepare("SELECT punkte.id, posten.kurz, posten.stamm, posten.kontakt, points, time FROM punkte LEFT JOIN posten ON punkte.von = posten.id WHERE an = '-1'")) {
				$stmt->execute();
				$stmt->store_result();
				$stmt->bind_result($id, $name, $stamm, $kontakt, $points, $time);
				while ($stmt->fetch()) {
					$posten++;
					echo '<li><a href="tel:'.$kontakt.'"><h2>Posten '.$name.' - '.$stamm.'</h2>
          <p>Punkte: '.$points.'<br>
          Zeitpunkt: '.$time.'</p></a>
          <a href="" class="punkte" id="'.$id.'"></a>
          </li>';
				}
        if ($posten == 0) echo "Kein Posten hat versucht eine Gruppe des anderen Laufs zu bewerten.";
      }
	?>
</ul>
<br />
<br />
<h3>Durchschnittsdifferenz</h3>
<ul data-role="listview" data-theme="a" data-inset="true">
	<?php
	$durchschnitt = array(0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0);
	$current = 0;
	if ($stmt = $mysqli->prepare("SELECT AVG(points) FROM posten, punkte, gruppen WHERE posten.id = von AND gruppen.id = an AND gruppen.stamm = posten.stamm GROUP BY von ORDER BY posten.kurz ASC")) {
		$stmt->execute();
		$stmt->store_result();
		$stmt->bind_result($avg);
		while ($stmt->fetch()) {
			$durchschnitt[$current] = $avg;
			$current++;
		}
	}
	$current = 0;
	$anzahl = 0;
	if ($stmt = $mysqli->prepare("SELECT posten.kurz, AVG(points) FROM posten, punkte, gruppen WHERE posten.id = von AND gruppen.id = an AND gruppen.stamm != posten.stamm GROUP BY von ORDER BY posten.kurz ASC")) {
		$stmt->execute();
		$stmt->store_result();
		$stmt->bind_result($kurz, $avg);
		while ($stmt->fetch()) {
			$durchschnitt[$current] = $durchschnitt[$current] - $avg;
			if($durchschnitt[$current] > 0) {
				echo '<li><h2>Posten '.$kurz.' bewertet '.$durchschnitt[$current].' Punkte über seinem Durchschnitt</h2></li>';
				$anzahl++;
			} else if ($durchschnitt[$current] < 0) {
				echo '<li><h2>Posten '.$kurz.' bewertet '.$durchschnitt[$current].' Punkte unter seinem Durchschnitt</h2></li>';
				$anzahl++;
			}
			$current++;
		}
		if ($anzahl == 0) {
			echo "Keine Auffälligkeiten.";
		}
	}
	?>
</ul>
<p align="center"><img src="files/wanze.png" style="max-width: 100%; max-height: 300px;"/></p>
<?php endif; ?>
