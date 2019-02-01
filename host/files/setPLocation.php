<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo '<fieldset class="ui-grid-a sideByside">';
  if ($stmt = $mysqli->prepare("SELECT id, kurz, x_axis, y_axis FROM posten ORDER BY kurz")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($id, $kurz, $xA, $yA);
    while ($stmt->fetch()) {
      echo '<div class="ui-block-a"><fieldset data-role="fieldcontain">
            <label for="x"><input data-mini="true" type="button" x="'.$xA.'" y="'.$yA.'" p="'.$kurz.'" class="map" value="'.$kurz.'"></label>
            <input type="text" class="coor" id="x" kurz="'.$kurz.'" value="'.$xA.'" placeholder="X">
            </fieldset></div>
            <div class="ui-block-b"><fieldset data-role="fieldcontain"><label for="y">&nbsp;</label>
            <input type="text" class="coor" id="y" kurz="'.$kurz.'" value="'.$yA.'" placeholder="Y">
            </fieldset>
            </div>';
    }
  } else echo "Etwas ist schiefgelaufen";
  echo '</fieldset>';
} else echo "Keine Berechtigung";
?>
