package main

import (
	"fmt"
	"net/http"
)

var scriptsFilter = `<link rel='stylesheet' type='text/css' href='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/css/theme.blue.min.css'>
<link rel="stylesheet" href="//code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css">
<script src='https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.0/jquery.min.js'></script>
<script src='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.9.1/jquery.tablesorter.min.js'></script>
<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js"></script>
<script type='text/javascript'>
  $(document).ready(function()
  {
    %s
    $('#myTable').tablesorter();
    if ($('#myTable2')) { $('#myTable2').tablesorter(); }
    $('#searchInput').keyup(function () {
    var data = this.value.split(' ');
    var jo = $('#myTable tbody').find('tr');
    if (this.value == '') { jo.show(); return; }
    jo.hide();
    jo.filter(function (i, v) {
      var $t = $(this);
      for (var d = 0; d < data.length; ++d) {
        if ($t.is(":contains('" + data[d] + "')")) { return true; } }
      return false;
    }).show();
    tChanged();
  }).focus(function () {
      this.value = '';
      $(this).css({ 'color': 'black' });
      $(this).unbind('focus');
  }).css({ 'color': '#C0C0C0' });
  tChanged();
 });
</script>
<script> function tChanged() { var x = $('#myTable tr:visible').length - 1; $('#rCount').text(x); }
</script>`

// PrintFilterHeader -- for sorted tables
func PrintFilterHeader(w http.ResponseWriter, header string, tableSort string, useDates bool) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	scr := fmt.Sprintf(scriptsFilter, "")
	fmt.Fprintf(w, "<html><head>%s</head><body>", scr)
	tSort := ""
	if len(tableSort) > 1 {
		tSort = "data-sortlist='" + tableSort + "'"
	}
	fmt.Fprintf(w, "<table id='myTable' class='tablesorter-blue' %s>", tSort)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><center><h1><b style='font-size:40px;color: #B9290A;text-shadow: 2px 2px 5px red;'>
	Welcome to the MSR Members directory</b></h1></center></body></html>`)
}

func Members(w http.ResponseWriter, r *http.Request) {
	PrintFilterHeader(w, "Current Members", "[[0,0]]", false)
	fmt.Fprintf(w, "<thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>",
		"ID", "Name", "NickName", "Email", "Address", "City", "Zip Code")
	memb := GetMembers()
	for _, o := range memb {
		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			o.id, o.name, o.nickName, o.email, o.address, o.city, o.zip)
	}
	fmt.Fprint(w, "</center></body></html>")

}
