<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php'; 

$wP;
$rP;

if ($stmt = $mysqli->prepare("SELECT sum(case when kategorie = '0' then 1 else 0 end) k0, sum(case when kategorie = '1' then 1 else 0 end) k1, sum(case when kategorie = '2' then 1 else 0 end) k2, sum(case when kategorie = '3' then 1 else 0 end) k3, sum(case when kategorie = '4' then 1 else 0 end) k4 FROM posten WHERE stufe = 'WöPo'")) {
	$stmt->execute();
	$stmt->store_result();
	$stmt->bind_result($wP[0], $wP[1], $wP[2], $wP[3], $wP[4]);
	while ($stmt->fetch()) {}
}
if ($stmt = $mysqli->prepare("SELECT sum(case when kategorie = '0' then 1 else 0 end) k0, sum(case when kategorie = '1' then 1 else 0 end) k1, sum(case when kategorie = '2' then 1 else 0 end) k2, sum(case when kategorie = '3' then 1 else 0 end) k3, sum(case when kategorie = '4' then 1 else 0 end) k4 FROM posten WHERE stufe = 'RoPo'")) {
	$stmt->execute();
	$stmt->store_result();
	$stmt->bind_result($rP[0], $rP[1], $rP[2], $rP[3], $rP[4]);
	while ($stmt->fetch()) {}
}
?>
