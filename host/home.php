<?php
include_once '../../includes/connect_gl.php';
require('../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) :
?>
<div class="ui-grid-b ui-responsive">
  <div class="ui-block-a"><input data-icon="back" value="Anmeldeseite" onclick="window.location.href='/stamm'" type="button"></div>
  <div class="ui-block-b"><input data-icon="home" value="Startseite" onclick="window.location.href='/'" type="button"></div>
  <div class="ui-block-c"><input data-icon="alert" value="Testseite" onclick="window.location.href='/host/test'" type="button"></div>
</div>
<h1>Hostseite <?php echo $login_session; ?></h1>
<div class="ui-grid-b ui-responsive">
<!--<div class="ui-block-a host" host="generatePosten" ><input data-icon="mail" data-theme="a" value="Alle Posten (Mail)" type="button"></div>
    <div class="ui-block-b host" host="generateGruppen"><input data-icon="mail" data-theme="a" value="Alle Gruppen (Mail)" type="button"></div>
--> <div class="ui-block-a host" host="gruppenGröße"><input data-icon="plus" data-theme="a" value="Gruppengröße" type="button"></div>
    <div class="ui-block-b host" host="gruppenKurz"><input data-icon="bullets" data-theme="a" value="Gruppenkürzel" type="button"></div>
    <div class="ui-block-c host" host="gruppenPunkte"><input data-icon="arrow-l" data-theme="a" value="Gruppen Punkte eintragen" type="button"></div>
    <div class="ui-block-a host" host="writeToPosten"><input data-icon="comment" data-theme="a" value="Posten Whatsapp" type="button"></div>
    <div class="ui-block-b host" host="writeLoginToPosten"><input data-icon="action" data-theme="a" value="Posten Whatsapp Zugangsdaten" type="button"></div>
    <div class="ui-block-c host" host="postenPunkte"><input data-icon="arrow-r" data-theme="a" value="Posten Punkte eintragen" type="button"></div>
    <div class="ui-block-a host" host="setPLocation"><input data-icon="location" data-theme="a" value="Posten Positionen" type="button"></div>
    <div class="ui-block-b host" host="editLogins"><input data-icon="check" data-theme="a" value="Login Passwörter" type="button"></div>
    <div class="ui-block-c host" host="antiCheat"><input data-icon="eye" data-theme="a" value="Wächterwanze" type="button"></div>
    <div class="ui-block-a host" host="startGruppen"><input data-icon="edit" data-theme="a" value="Posten Startgruppen" type="button"></div>
    <div class="ui-block-b host" host="postenDurch"><input data-icon="plus" data-theme="b" value="Postendurchschnitt" type="button"></div>
    <div class="ui-block-c host" host="gruppenDurch"><input data-icon="plus" data-theme="b" value="Gruppendurchschnitt" type="button"></div>
    <div class="ui-block-b host" host="postenTop"><input data-icon="star" data-theme="b" value="Postentabelle" type="button"></div>
    <div class="ui-block-c host" host="gruppenTop"><input data-icon="star" data-theme="b" value="Gruppentabelle" type="button"></div>
    <div class="ui-block-a host" host="fill0"><input data-icon="edit" data-theme="b" value="Fülle Punkte" type="button"></div>
    <div class="ui-block-b host" host="deletePoints"><input data-icon="forbidden" data-theme="b" value="Lösche ALLE Punkte" type="button"></div>
    <div class="ui-block-c host" host="backupPunkte"><input data-icon="arrow-d" data-theme="b" value="Punkte Backup erstellen" type="button"></div>
    <div class="ui-block-a host" host="query"><input data-icon="grid" data-theme="b" value="Query" type="button"></div>
</div>
<br>
<br>
<div id="content">
</div>
<?php endif; ?>
