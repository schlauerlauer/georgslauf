<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<h2>Gruppenkürzel bearbeiten</h2>
<ul data-role="listview" data-count-theme="b" data-inset="true" data-autodividers="true" data-filter="true" data-filter-placeholder="Gruppen Filtern">
	<?php
    $gruppen = 0;
		if ($stmt = $mysqli->prepare("SELECT id, kurz, name, stufe, stamm FROM gruppen ORDER BY stufe, stamm, name asc")) {
			//$stmt->bind_param('s', $user);
			$stmt->execute();
			$stmt->store_result();
			$stmt->bind_result($id, $kurz, $name, $stufe, $stamm);
			while ($stmt->fetch()) {
        $gruppen++;
				echo '<li><a href="#" class="gkurz" id="'.$id.'" kurz="'.$kurz.'">'.$Stufe[$stufe].' &nbsp;'.$stamm.' - '.$name.'<span class="ui-li-count">'.$kurz.'</span></a></li>';
			}
		}
    if($gruppen == 0) echo "Noch keine Gruppen angemeldet";
	?>
</ul>
<?php endif;
?>
