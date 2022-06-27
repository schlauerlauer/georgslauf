<?php

include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<h2>Alle Posten</h2>
<div data-role="collapsible-set" data-filter="true" data-filter-placeholder="Posten Filtern" data-theme="a">
	<?php
		$pval = 0;
		$totalS = 0;
		$totalV = 0;
		if ($stmt = $mysqli->prepare("SELECT id, color, kurz, name, stamm, kategorie, beschreibung, kontakt, anzahl, veggie, material, ort, sonstiges, stufe FROM posten ORDER BY stufe DESC, stamm, name")) {
			$stmt->execute();
			$stmt->store_result();
			$stmt->bind_result($id, $color, $kurz, $name, $stamm, $kat, $desc, $kont, $val, $veg, $mat, $ort, $sonst, $stufe);
			while ($stmt->fetch()) {
				$totalS += $val;
				$totalV += $veg;
				echo '<div data-role="collapsible"><h3>'.$kurz.' '.$stufe.' - '.$stamm.' - '.$Kat[$kat].'</h3><p>
				<strong>'.$name.'</strong><br>"'.$desc.'"<br><br>Leiter '.$val.' ('.$veg.' Veggie)<br>Kontakt <a href="tel:'.$kont.'">'.$kont.'</a>';
				if (!empty($mat)) echo '<br>Material '.$mat;
				if (!empty($ort)) echo '<br>Ort '.$ort;
				if (!empty($sonst)) echo '<br>Sonstiges '.$sonst;
				echo '
					<div class="ui-grid-a">
						<div class="ui-block-a color" id="'.$id.';">
							<div class="ui-bar ui-bar-a" style="height:50px; background-color:'.$Farben[$color].';"></div>
						</div>
						<div class="ui-block-b">
							<div class="ui-bar ui-bar-b" border="0" style="height:50px; background-color:white;">
								<div class="ui-field-contain">
									<label for="'.$id.'" style="color:black;">Postenkürzel</label>
									<input id="'.$id.';" class="kurz" placeholder="Postenkürzel" maxlength="1" data-mini="true" value="'.$kurz.'"/>
								</div>
							</div>
						</div>
					</div></p></div>';
				$pval++;
			}
			if ($pval == 0) echo "Noch keine Posten angemeldet";
		}
	?>
</div>
<br>
<h3>Insgesamt <?php echo $pval; ?> Posten mit <?php echo $totalS; ?> Leitern, davon <?php echo $totalV; ?> Vegetarier</h3>
<?php endif; ?>
