<?php
include('session/login.php');
include('host/settings.php');

if(isset($_SESSION['login_user'])) : ?>
	<h3> Als <?php echo $Rollen[$_SESSION['rolle']].' '.$_SESSION['login_user']; ?> angemeldet</h3>
	<?php if($_SESSION['rolle'] == 2) : ?>
	<button class="ui-btn ui-icon-carat-r ui-btn-icon-left" onclick="window.location.href='/stamm'">Zur Stammseite</button>
	<?php elseif ($_SESSION['rolle'] == 1) : ?>
	<button class="ui-btn ui-icon-carat-r ui-btn-icon-left" onclick="window.location.href='/posten'">Zur Postenseite</button>
<?php elseif ($_SESSION['rolle'] >= 3) : ?>
	<button class="ui-btn ui-icon-carat-r ui-btn-icon-left" onclick="window.location.href='/stamm'">Zur Stammseite</button>
	<button class="ui-btn ui-btn-b ui-icon-gear ui-btn-icon-left" onclick="window.location.href='/host'">Zur Hostseite</button>
	<?php endif; ?>
	<br>
	<button class="ui-btn ui-btn-b ui-icon-user ui-btn-icon-left" onclick="window.location.href='/session/logout.php'">Abmelden</button>
<?php else : ?>
	<h2>Login</h2>
	<p>Posten Login<br><br>
	<!--	Achtung: Änderungen an Posten und Laufgruppen nur noch per E-Mail!<br><br> -->
	</p>
	<br>
	<label>Benutzername</label>
	<input id="username" name="username" placeholder="Benutzername" type="text">
	<label>Passwort</label>
	<input id="password" name="password" placeholder="**********" type="password">
	<button class="ui-btn" onclick="login()">Login</button>
<?php endif; ?>
<br>
<br>
<p align="center"><img src="../res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
