<?php
require('../session/session.php');
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
include_once '../host/settings.php';
if(isset($login_session) && $_SESSION['rolle'] == 1) : ?>
<h1>Gruppenbewertungen</h1>
<button class="ui-btn ui-btn-inline ui-icon-home ui-btn-icon-left" onclick="window.location.href='/'">Startseite</button>
<h2>Posten <?php echo $login_session;?></h2>
<ul data-role="listview" data-theme="a" data-inset="true">
	<?php
    if($Punkte[0] == true) {

			$user = $login_session;
			$alter = 0;
      if($user >= $PostenStufe) $alter = 2;

			$gruppen = 0;
			if ($stmt = $mysqli->prepare("SELECT g.id, g.kurz, g.name, g.size, g.stufe, g.stamm, punkte.points FROM gruppen g LEFT JOIN punkte ON g.id = punkte.an
			WHERE (von = (SELECT id from posten where kurz = ?) OR von IS NULL)
			 ORDER BY kurz ASC")) {
				$stmt->bind_param('s', $user);
				$stmt->execute();
				$stmt->store_result();
				$stmt->bind_result($id, $kurz, $name, $size, $stu, $sta, $points);
				while ($stmt->fetch()) {
					$gruppen++;
					echo '<li><h2>'.$kurz.' '.$name.'</h2>
          <p>'.$size.' '.$Stufe[$stu].' - '.$sta.'</p>
          <p class="ui-li-aside"><input id="'.$id.'" class="punkte" style="max-width:60px;" data-mini="true" type="number" min="0" max="100" value="'.$points.'"/></p>
          </li>';
				}
        if ($gruppen == 0) echo "Noch keine Gruppen angemeldet";
      }
		} else echo '<h3 class="ui-bar ui-bar-a ui-corner-all" align="center">'.$Punkte[1].'</h3>';
	?>
</ul>
<br>
<br>
<br>
<p align="center"><img src="/res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
<?php endif;
/*
Besser aber einträge fehlen, wenn von anderem posten eingetragen wurde -> alles mit nullen füllen
SELECT g.id, g.kurz, g.name, g.stufe, g.stamm, punkte.points FROM gruppen g LEFT JOIN punkte ON g.id = punkte.an
WHERE (von = (SELECT id from posten where kurz = 'B') OR von IS NULL)
AND (g.stufe = 0 OR g.stufe = 1) ORDER BY kurz ASC
*/
?>
