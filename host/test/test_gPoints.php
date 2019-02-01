<?php
include_once '../../../includes/connect_gl.php';
include_once '../settings.php';
//include_once '/session/session.php';
?>
<h1>Gruppen Punkte</h1>
<h2>von Posten wX</h2>
<ul data-role="listview" data-theme="a" data-inset="true">
	<?php
    if($Punkte[0] == true) {
    //TODO Tabindex
			$user = 'wX';
      $pID = 114; //beim login abfragen und setzen!
      $alter = 0; //Altersstufe des Postens
      if(substr($user,0,1) == "r") $alter = 2;
			$gruppen = 0;
			if ($stmt = $mysqli->prepare("SELECT g.id, g.name, g.stufe, g.stamm, punkte.points FROM gruppen g LEFT JOIN punkte ON g.id = punkte.an
        WHERE von = 114 OR von IS NULL AND g.stufe = 0 OR g.stufe = 1 ORDER BY id asc")) {
				//$stmt->bind_param('sss', '114', '0', '1'); //funktioniert nicht
				$stmt->execute();
				$stmt->store_result();
				$stmt->bind_result($id, $name, $stu, $sta, $points);
				while ($stmt->fetch()) {
					$gruppen++; //KÜRZEL! unten einfügen
					echo '<li><h2>'.$name.'</h2>
          <p>'.$Stufe[$stu].' - '.$sta.'</p>
          <p class="ui-li-aside"><input id="'.$id.'" class="punkte" data-mini="true" type="number" min="0" max="100" value="'.$points.'"/></p>
          </li>';
				}
        if ($gruppen == 0) echo "Noch keine Gruppen angemeldet";
      }
		} else echo $Punkte[1];
	?>
</ul>
<br>
<br>
<br>
<p align="center"><img src="/res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
