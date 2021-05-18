<?php

include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');


if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<h2>Alle Gruppen</h2>
<ul data-role="listview" data-count-theme="b" data-inset="true" data-autodividers="true" data-filter="true" data-filter-placeholder="Gruppen Filtern">
	<?php
		$gruppen = 0;
		$totalS = 0;
		$totalV = 0;
		if ($stmt = $mysqli->prepare("SELECT id, name, stufe, stamm, size, veggies FROM gruppen ORDER BY stamm, stufe, name asc")) {
			$stmt->execute();
			$stmt->store_result();
			$stmt->bind_result($id, $name, $stufe, $stamm, $size, $veggies);
			while ($stmt->fetch()) {
				$gruppen++;
				$totalS += $size;
				$totalV += $veggies;
				echo '<li><a href="#">'.$stamm.'  '.$Stufe[$stufe].' - '.$name.'<span class="ui-li-count">'.$size.'</span></a></li>';
				//<a href="#" class="ui-btn ui-btn-inline ui-icon-delete ui-btn-icon-notext delete" type="g" name='.$name.' id='.$id.'></a>
				//</div>';
			}
			if ($gruppen == 0) echo "Noch keine Gruppen angemeldet";
		}
	?>
</ul>
<br>
<h4>Insgesamt <?php echo $gruppen; ?> Laufgruppen mit <?php echo $totalS; ?> Teilnehmern, davon <?php echo $totalV; ?> Vegetarier</h4>
<?php endif;
//SELECT an, name, size, stufe, stamm, sum(points) FROM gruppen, punkte WHERE an = name_varchar GROUP BY an ORDER BY points DESC
?>
