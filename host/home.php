<?php
include_once '../../includes/connect_gl.php';
require('../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<h1>Georgslauf 2020 <?php echo $login_session; ?></h1>
<div class="ui-grid-b ui-responsive">
  <div class="ui-block-a"><input data-icon="back" value="Anmeldeseite" onclick="window.location.href='/stamm'" type="button"></div>
  <div class="ui-block-b"><input data-icon="home" value="Startseite" onclick="window.location.href='/'" type="button"></div>
</div>
<h3>Posten & Gruppen</h3>
<div class="ui-grid-b ui-responsive">
  <div class="ui-block-a host" host="gruppen"><input data-icon="carat-r" data-theme="a" value="Alle Gruppen" type="button"></div>
  <div class="ui-block-b host" host="posten"><input data-icon="carat-r" data-theme="a" value="Alle Posten" type="button"></div>
  <div class="ui-block-c host" host="gruppenKurz"><input data-icon="bullets" data-theme="a" value="Gruppenkürzel" type="button"></div>
  <div class="ui-block-a host" host="setPLocation"><input data-icon="location" data-theme="a" value="Posten Positionen" type="button"></div>
  <div class="ui-block-b host" host="startGruppen"><input data-icon="edit" data-theme="a" value="Posten Startgruppen" type="button"></div>
  <div class="ui-block-c " host="deleteAll"><input disabled data-icon="edit" data-theme="a" value="Alle Posten und Gruppen löschen" type="button"></div>
</div>
<h3>Punkte</h3>
<div class="ui-grid-b ui-responsive">
  <div class="ui-block-a host" host="gruppenPunkte"><input data-icon="arrow-l" data-theme="a" value="Gruppen Punkte eintragen" type="button"></div>
  <div class="ui-block-b host" host="postenPunkte"><input data-icon="arrow-r" data-theme="a" value="Posten Punkte eintragen" type="button"></div>
  <div class="ui-block-c host" host="antiCheat"><input data-icon="eye" data-theme="a" value="Wächterwanze" type="button"></div>
  <div class="ui-block-a host" host="query"><input data-icon="grid" data-theme="b" value="Query" type="button"></div>
  <div class="ui-block-b " host="fill0"><input data-icon="edit" data-theme="b" disabled style="background:red;" value="Alle Punkte zurücksetzen" type="button"></div>
  <div class="ui-block-c " host="deletePoints"><input data-icon="forbidden" disabled style="background:red;" data-theme="b" value="Alle Punkte löschen" type="button"></div>
  <!--<div class="ui-block-a host" host="backupPunkte"><input data-icon="arrow-d" disabled data-theme="b" value="Punkte Backup erstellen" type="button"></div>-->
  <!--<div class="ui-block-c host" host="gruppenGröße"><input data-icon="plus" disabled data-theme="a" value="Gruppengröße" type="button"></div>-->
</div>
<h3>Kommunikation</h3>
<div class="ui-grid-b ui-responsive">
  <div class="ui-block-a host" host="writeToPosten"><input data-icon="comment" data-theme="a" value="Posten Whatsapp" type="button"></div>
  <div class="ui-block-b host" host="posten"><input data-icon="phone" data-theme="a" value="Posten anrufen" type="button"></div>
  <div class="ui-block-c host" host="writeLoginToPosten"><input data-icon="action" data-theme="a" value="Posten Whatsapp Zugangsdaten" type="button"></div>
  <div class="ui-block-a host" host="editLogins"><input data-icon="check" data-theme="a" value="Login Passwörter" type="button"></div>
  <div class="ui-block-b host" host="generatePosten" ><input data-icon="mail" data-theme="a" value="Alle Posten (Mailformat)" type="button"></div>
  <div class="ui-block-c host" host="generateGruppen"><input data-icon="mail" data-theme="a" value="Alle Gruppen (Mailformat)" type="button"></div>
</div>
<h3>Siegerehrung</h3>
<div class="ui-grid-b ui-responsive">
  <div class="ui-block-a host" host="postenTop"><input data-icon="star" data-theme="b" value="Posten Sieger" type="button"></div>
  <div class="ui-block-b host" host="gruppenTop"><input data-icon="star" data-theme="b" value="Gruppen Sieger" type="button"></div>
  <div class="ui-block-c host" host="presentation"><input data-icon="arrow-d" data-theme="b" value="Erstelle Präsentation" type="button"></div>
  <div class="ui-block-a host" host="postenDurch"><input data-icon="plus" data-theme="b" value="Postendurchschnitt" type="button"></div>
  <div class="ui-block-b host" host="gruppenDurch"><input data-icon="plus" data-theme="b" value="Gruppendurchschnitt" type="button"></div>
</div>
<br>
<br>
<div id="content">
</div>
<br>
<br>
<div data-role="footer" data-tap-toggle="false" style="overflow:hidden; margin: 0px;">
		<h2><img src="/res/logo.png" style="max-height:150px;"/></h2>
</div>
<?php endif; ?>
