package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var styles = `<link rel='stylesheet' type='text/css' href='/static/index.css'>`

var MenuHdr = `
<div class="dropdown" style='margin-top: 14px;margin-left: 12px;'>
  <b style='vertical-align: bottom;font-size:30px;color: #B9290A;text-shadow: 2px 2px 5px red;'>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;%s</b>
</div>
<br><p></p>`

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

var scripts = `<link rel='stylesheet' type='text/css' href='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/css/theme.blue.min.css'>
        <script src='https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.0/jquery.min.js'></script>
        <script src='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.9.1/jquery.tablesorter.min.js'></script>
        <script>
        $(document).ready(function() {
  $('#myTable').tablesorter(); 
  $('#myTable2').tablesorter(); } ); </script>`

// PrintFilterHeader -- for sorted tables
func PrintFilterHeader(w http.ResponseWriter, header string, tableSort string, useDates bool) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	scr := fmt.Sprintf(scriptsFilter, "")
	fmt.Fprintf(w, "<html><head>%s</head><body>", scr)
	fmt.Fprintf(w, MenuHdr, header)
	tSort := ""
	if len(tableSort) > 1 {
		tSort = "data-sortlist='" + tableSort + "'"
	}
	fmt.Fprintf(w, "<table id='myTable' class='tablesorter-blue' %s>", tSort)
}

// PrintUIHeader -- Basic UI
func PrintUIHeader(w http.ResponseWriter, header string) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html><head>%s</head><body>", styles+scripts)
	fmt.Fprintf(w, MenuHdr, header)
	fmt.Fprintf(w, "<br />")
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><center><h1><b style='font-size:40px;color: #B9290A;text-shadow: 2px 2px 5px red;'>
	Welcome to the MSR Members directory</b></h1></center></body></html>`)
}

func Members(w http.ResponseWriter, r *http.Request) {
	PrintFilterHeader(w, "MSR Current Members", "[[0,0]]", false)
	fmt.Fprintf(w, "<thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>",
		"ID", "Name", "NickName", "Email", "Address", "City", "Zip Code")
	memb := GetMembers()
	for _, o := range memb {
		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>",
			o.id, o.name, o.nickName, o.email, o.address, o.city, o.zip)
	}
	fmt.Fprint(w, "</center></body></html>")

}

func MembersS(w http.ResponseWriter, r *http.Request) {
	PrintFilterHeader(w, "MSR Current Members", "[[0,0]]", false)
	fmt.Fprintf(w, "<thead><tr><th>%s</th><th>%s</th><th>%s</th></thead><tbody>",
		"Name", "NickName", "Email")
	memb := GetMembersS()
	for _, o := range memb {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td></tr>",
			o.name, o.nickName, o.email)
	}
	fmt.Fprint(w, "</center></body></html>")

}

// EditMember -- Edit member
func EditMember(w http.ResponseWriter, r *http.Request) {
	var mb Member
	vars := mux.Vars(r)
	fmt.Printf("EditCust, %v\n", r.Form)
	memberID, _ := strconv.Atoi(vars["memberId"])
	if memberID > 0 {
		mb = GetMember(memberID)
		fmt.Printf("Edit Member, ID=%d, Name=%s\n", memberID, mb.name)
	} else {
		fmt.Println("Add Member")
		mb.id = 0
		mb.name = "New"
	}
	//tmpl := template.Must(template.ParseFiles("edCustomer.html"))
	//TMPLCus.Execute(w, cust)
	TMPLAll.ExecuteTemplate(w, "edMember.html", mb)
}

// SaveMember -- Save Member
func SaveMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mID, _ := strconv.Atoi(vars["memberId"])
	r.ParseForm()
	fmt.Printf("Saving Member, ID=%d, Name=%s\n", mID, r.FormValue("iName"))
	//log.Println(r.Form)
	var m Member
	m.id = mID
	m.name = r.FormValue("iName")
	m.nickName = r.FormValue("iNickName")
	m.address = r.FormValue("iAdress")
	m.zip, _ = strconv.Atoi(r.FormValue("iZip"))
	m.city = r.FormValue("iCity")
	m.email = r.FormValue("iEmail")
	m.notes = r.FormValue("iNotes")
	id := UpdMember(m)
	if mID == 0 {
		PrintUIHeader(w, "Saving Member")
		fmt.Fprintf(w, "<p>New Member saved, Name=%s --  Go to <a href='/members'>Search</a>.</p>",
			m.name)
	} else {
		log.Printf("Data saved for the Customer, ID=%d, Name=%s", id, m.name)
		http.Redirect(w, r, "/members", http.StatusFound)
	}
}
