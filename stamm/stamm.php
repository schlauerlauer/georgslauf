<?php
include_once '../host/settings.php';
include_once 'pGet.php';
include_once '../session/session.php';
?>
<button class="ui-btn ui-btn-inline ui-icon-home ui-btn-icon-left" onclick="window.location.href='/'">Zurück zur Startseite</button>
<?php if ($login_session == $Host) : ?>
<button class="ui-btn ui-btn-b ui-btn-inline ui-icon-gear ui-btn-icon-left" onclick="window.location.href='/host'">Hostseite</button>
<?php endif; ?>
<?php if($_SESSION['rolle'] >= 2) : ?>
<h1>Anmeldung Stamm <?php echo $login_session; ?></h1>
Dokumente
<!--<a target="_blank" class="ui-btn ui-btn-b ui-icon-alert ui-btn-icon-left" href="/res/Informationen.pdf">Informationen anzeigen (.pdf)</a>-->
<a target="_blank" class="ui-btn ui-icon-edit ui-btn-icon-left" href="/res/Interne_Anmeldung.docx">Vorlage für eure interne Anmeldung (.docx)</a>
bei Fragen: <a href="mailto:gl20@stamm-prm.de">gl20@stamm-prm.de</a> oder <a href="tel:01709462008">anrufen (01709462008)</a> / <a target="_blank" href="https://api.whatsapp.com/send?phone=491709462008&text=Frage%20von%20Stamm%20<?php echo $login_session; ?>%3A%0D%0A%0D%0A">WhatsApp (01709462008)</a>
<br>
<br>
<br>
<?php if ($Anmeldung == true) : ?>
<h2>Gruppe anmelden</h2>
<form id="g_form">
    <fieldset data-role="controlgroup">
       <legend>Stufe<span style="color:red;">*</span>&nbsp;&nbsp;
	   <a href="#" class="help ui-btn ui-btn-inline ui-icon-info ui-btn-icon-notext ui-corner-all" help="Gruppenstufen"></a>
	   </legend>
        <input name="stufe" id="s_1" value="0" type="radio">
        <label for="s_1">Wölflinge</label>
        <input name="stufe" id="s_2" value="1" type="radio">
        <label for="s_2">Jupfis</label>
        <input name="stufe" id="s_3" value="2" type="radio">
        <label for="s_3">Pfadis</label>
		<input name="stufe" id="s_4" value="3" type="radio">
        <label for="s_4">Rover</label>
    </fieldset>
	<label for="g_name">Gruppenname<span style="color:red;">*</span>&nbsp;&nbsp;
	<a href="#" class="help ui-btn ui-btn-inline ui-icon-info ui-btn-icon-notext ui-corner-all" help="Gruppennamen"></a>
	</label>
    <input id="g_name" value="" type="text" maxlength="50">
	<label for="g_anzahl">Anzahl Kinder 5-12<span style="color:red;">*</span>&nbsp;&nbsp;
	<a href="#" class="help ui-btn ui-btn-inline ui-icon-info ui-btn-icon-notext ui-corner-all" help="Gruppengrößen"></a></label>
    <input data-clear-btn="false" id="g_anzahl" value="5" type="number" min="4" max="15">
    <label for="g_veggie">Davon Vegetarier</label>
      <input data-clear-btn="false" id="g_veggie" value="0" type="number" min="0" max="15">
	<br>
	<span style="color:red;">*</span> benötigte Angaben
	<a href="" class="ui-btn ui-btn-b ui-icon-check ui-btn-icon-left save" id="g_save">Gruppe speichern</a>
</form>
<br>
<br>
<h2>Posten anmelden</h2>
<form id="p_form">
	<fieldset data-role="controlgroup">
       <legend>Stufe<span style="color:red;">*</span>&nbsp;&nbsp;
	   <a href="#" class="help ui-btn ui-btn-inline ui-icon-info ui-btn-icon-notext ui-corner-all" help="Stufen"></a></legend>
        <input name="p_stufe" id="p_1" value="WöPo" type="radio">
        <label for="p_1">Wölfling / Jupfi (WöPo)</label>
        <input name="p_stufe" id="p_2" value="RoPo" type="radio">
        <label for="p_2">Pfadi / Rover (RoPo)</label>
    </fieldset>
</form>
<form id="k_form">
	<fieldset data-role="controlgroup">
       <legend>Kategorie<span style="color:red;">*</span>&nbsp;&nbsp;
	   <a href="#" class="help ui-btn ui-btn-inline ui-icon-info ui-btn-icon-notext ui-corner-all" help="Kategorien"></a></legend>
		<?php
			for($i = 0; $i < sizeof($Kat); $i++) {
				echo '<input name="kategorie" id="k_'.$i.'" value="'.$i.'" type="radio">
				<label for="k_'.$i.'">'.$Kat[$i].'</label>';
			}
		?>
    </fieldset>
	<label for="p_name">Postenname<span style="color:red;">*</span></label>
    <input id="p_name" value="" type="text" maxlength="50">
	<label for="p_desc">Postenbeschreibung<span style="color:red;">*</span></label>
    <input id="p_desc" value="" type="text" placeholder="Was wird an dem Posten gemacht?" maxlength="255">
	<label for="p_anzahl">Anzahl Leiter<span style="color:red;">*</span></label>
	<input data-clear-btn="false" id="p_anzahl" value="2" type="number" min="1" max="20">
	<label for="p_veggie">Davon Vegetarier<span style="color:red;">*</span></label>
	<input data-clear-btn="false" id="p_veggie" value="0" type="number" min="0" max="20">
	<label for="p_kont">Handynummer eines Postenleiters<span style="color:red;">*</span></label>
    <input id="p_kont" value="" type="text" placeholder="Wichtig für den Startschuss (Whatsapp)" maxlength="100">
	<label for="p_mat">Benötigtes Material</label>
    <input id="p_mat" value="" type="text" placeholder="Bierbänke ..." maxlength="200">
	<label for="p_ort">Bevorzugter Ort</label>
    <input id="p_ort" value="" type="text" placeholder="Wiese / Bürgersteig / Nähe Parkplatz ..." maxlength="200">
	<label for="p_sonst">Sonstiges</label>
    <input id="p_sonst" value="" type="text" placeholder="Sonstige Wünsche / Anmerkungen ..." maxlength="200">
	<br>
	<span style="color:red;">*</span> benötigte Angaben
	<a href="" class="ui-btn ui-btn-b ui-icon-check ui-btn-icon-left save" id="p_save">Posten speichern</a>
</form>
<?php else : ?>
<h2 align="center">Die Anmeldung ist nur noch per Email möglich</h2>
<?php endif; ?>
<?php else : ?>
<h2>keine Rechte</h2>
<?php endif; ?>
<br>
<br>
<div align="center"><a target="_blank" href="http://pfadi-fc.de"><img src="../res/fc.png"/></a></div>
