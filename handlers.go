package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	styles = `<link rel='stylesheet' type='text/css' href='/static/index.css'>`

	MenuHdr = `
<div class="dropdown" style='margin-top: 14px;margin-left: 12px;'>
  <button class="dropbtn"><img src="/static/Notebook.jpg" style="height:32px;">&nbsp;&nbsp;Menu...</button>
  <div class="dropdown-content">
  <a href="/">Home</a>
  <a href="/members">Members</a>
  </div>
  <b style='vertical-align: bottom;font-size:30px;color: #B9290A;text-shadow: 2px 2px 5px red;'>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;%s</b>
</div>
<br><p></p>`

	scriptsFilter = `<link rel='stylesheet' type='text/css' href='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/css/theme.blue.min.css'>
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

	scripts = `<link rel='stylesheet' type='text/css' href='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/css/theme.blue.min.css'>
        <script src='https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.0/jquery.min.js'></script>
        <script src='https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.9.1/jquery.tablesorter.min.js'></script>
        <script>
        $(document).ready(function() {
  $('#myTable').tablesorter(); 
  $('#myTable2').tablesorter(); } ); </script>`
)

// PrintFilterHeader -- for sorted tables
func PrintFilterHeader(w http.ResponseWriter, header string, tableSort string, useDates bool) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	scr := fmt.Sprintf(scriptsFilter, "")
	fmt.Fprintf(w, "<html><head>%s</head><body>", styles+scr)
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

// GetMemberHref - get the link to edit staff member
func GetMemberHref(id int) string {
	var ref = fmt.Sprintf("<a href='/editmember/%d'>Edit %d</a>", id, id)
	return ref
}

// Login -- Login page
func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("Username")
	password := r.FormValue("Password")
	fmt.Println("Username=" + username + ", password=" + password)
	if len(username) > 0 && len(password) > 0 {
		pwd := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		userAccess := LoginUser(username, pwd)
		if userAccess > 0 {
			expiration := time.Now().Add(20 * time.Hour)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
			http.SetCookie(w, &cookie)
			acookie := http.Cookie{Name: "uaccess", Value: strconv.Itoa(userAccess), Expires: expiration}
			http.SetCookie(w, &acookie)
			http.Redirect(w, r, "/", http.StatusFound)
			//ActiveUsers[username] = userAccess
		}
	}
	fmt.Fprint(w, "<html><body><center><h1>Not authenticated</h1><p><a href='/static/Login.html'>Login page...</a></p></body></html>")
}

// IsAdminUser -- Check user cookie to see if he has admin rights
func IsAdminUser(r *http.Request) bool {
	cookie, _ := r.Cookie("uaccess")
	uacc := 0
	if cookie != nil {
		uacc, _ = strconv.Atoi(cookie.Value)
		fmt.Printf("User Cookie, need a 9, cookie has %d\n", uacc)
	}
	return uacc == 9
}

func NotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/notadminuser", http.StatusFound)
}

func NotAdminUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><center><h1><b style='font-size:40px;color: #B9290A;text-shadow: 2px 2px 5px red;'>
	MSR Members</b></h1>
	<p>Not allowed</p><p>You need to be Admin user to get access</p></body></html>`)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><center><h1><b style='font-size:40px;color: #B9290A;text-shadow: 2px 2px 5px red;'>
	Welcome to the MSR Members directory</b></h1></body></html>`)
	//fmt.Fprint(w, `<p><img src="/static/Notebook.jpg" alt="Members" style="height:240px;">`)
	fmt.Fprint(w, `<p><img src="/static/MSR.JPEG" alt="MSR" style="height:325px;">`)
	fmt.Fprint(w, "<p><a href='/members'>Manage Members</a></p>")
}

func Members(w http.ResponseWriter, r *http.Request) {
	PrintFilterHeader(w, "MSR Current Members", "[[1,0]]", false)
	fmt.Fprintf(w, "&nbsp;&nbsp;&nbsp;&nbsp;")
	fmt.Fprintf(w, "Click here to <a href='/editmember/0'>Add new Member</a></center></div>")
	fmt.Fprintf(w, "<thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>",
		"ID", "Name", "Rank", "Year", "Phone", "Address", "City", "State", "Zip Code", "Email", "Birth Day", "Status")
	memb := GetMembers()
	for _, o := range memb {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%d</td></tr>",
			GetMemberHref(o.ID), o.Name, o.Rank, o.Since, o.Phone, o.Address, o.City, o.State, o.Zip, o.Email, o.BirthDay, o.Status)
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
			o.Name, o.Email, o.Phone)
	}
	fmt.Fprint(w, "</center></body></html>")

}

// EditMember -- Edit member
func EditMember(w http.ResponseWriter, r *http.Request) {
	var mb Member
	vars := mux.Vars(r)
	fmt.Printf("EditMember, %v\n", r.Form)
	if IsAdminUser(r) {
		memberID, _ := strconv.Atoi(vars["memberId"])
		fmt.Printf("Member.ID, %d\n", memberID)
		if memberID > 0 {
			mb = GetMember(memberID)
			fmt.Printf("Edit Member, ID=%d, Name=%s\n", memberID, mb.Name)
		} else {
			fmt.Println("Add Member")
			mb.ID = 0
			mb.Name = "New"
			mb.Address = ""
			mb.Email = "test@test.com"
			mb.Status = 0
			mb.Zip = 0
		}
		tmpl := template.Must(template.New("test").Parse(`
	<html>
    <link rel='stylesheet' type='text/css' href='/static/index.css'>
    <script src='https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js'></script>
    <body><form action="/savemember/{{.ID}}" target="_self" method="POST"><table>
    <tr><td></td>
                    <td>{{if (gt .ID 0)}}
                            <b style='font-size:30px;color: #B9290A;text-shadow: 2px 2px 5px red;'>Edit a Member</b>
                        {{else}}
                            <b style='font-size:30px;color: #B9290A;text-shadow: 2px 2px 5px red;'>Add a new Member</b>
                        {{end}}</td>
                </tr><tr>
				<td><label for="iName">Name:</label></td>
				<td><input type="text" id="iName" name="iName" size="64" value="{{.Name}}"/></td>
			</tr><tr>
				<td><label for="iRank">Rank:</label></td>
				<td><input type="number" id="iRank" name="iRank" size="4" value="{{.Rank}}"/>
				<label for="iSince">&nbsp;&nbsp;Year:</label>
				<input type="text" id="iSince" name="iSince" size="10" value="{{.Since}}"/></td>
			</tr><tr>
				<td><label for="iPhone">Phone:</label></td>
				<td><input type="text" id="iPhone" name="iPhone" size="12" value="{{.Phone}}"/></td>
			</tr><tr>
				<td><label for="iAddress">Address:</label></td>
				<td><input type="text" id="iAddress" name="iAddress" size="64" value="{{.Address}}"/></td>
			</tr><tr>
				<td><label for="iCity">City:</label></td>
				<td><input type="text" id="iCity" name="iCity" size="50" value="{{.City}}"/></td>
			</tr><tr>
				<td><label for="iState">State:</label></td>
				<td><input type="text" id="iState" name="iState" size="3" value="{{.State}}"/></td>
			</tr><tr>
				<td><label for="iZip">Zip Code:</label></td>
				<td><input type="number" id="iZip" name="iZip" size="8" value="{{.Zip}}"/></td>
			</tr><tr>
				<td><label for="iEmail">Email:</label></td>
				<td><input type="text" id="iEmail" name="iEmail" size="50" value="{{.Email}}"/></td>
			</tr><tr>
				<td><label for="iBirthday">Birthday:</label></td>
				<td><input type="text" id="iBirthday" name="iBirthday" size="10" value="{{.BirthDay}}"/>
				<label for="iStatus">&nbsp;&nbsp;Status:</label>
				<input type="number" id="iStatus" name="iStatus" size="2" value="{{.Status}}"/>
				<label>&nbsp;&nbsp;0=Member, 1=Prospect, 2=On Leave</label></td>
			</tr><tr>
				<td><label for="iNotes">Notes:</label></td>
				<td><textarea id="iNotes" name="iNotes" rows="4" cols=50>{{printf "%s" .Notes}}</textarea></td>
			</tr><tr><td> </td></tr><tr>
				<td><input type="submit" class="gbutton" value="Save"></td>
				<td><a id="iBack" class="gbutton" href='javascript:history.back()'>Cancel</a></td></tr>
			</table>
        </form></body></html>`))
		fmt.Printf("Execute Template for %s\n", mb.Name)
		tmpl.Execute(w, mb)
	} else {
		NotAllowed(w, r)
	}
}

// SaveMember -- Save Member
func SaveMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if IsAdminUser(r) {
		mID, _ := strconv.Atoi(vars["memberId"])
		r.ParseForm()
		fmt.Printf("Saving Member, ID=%d, Name=%s\n", mID, r.FormValue("iName"))
		//log.Println(r.Form)
		var m Member
		m.ID = mID
		m.Rank, _ = strconv.Atoi(r.FormValue("iRank"))
		m.Since = r.FormValue("iSince")
		m.Name = r.FormValue("iName")
		m.Phone = r.FormValue("iPhone")
		m.Address = r.FormValue("iAddress")
		m.Zip, _ = strconv.Atoi(r.FormValue("iZip"))
		m.City = r.FormValue("iCity")
		m.State = r.FormValue("iState")
		m.Email = r.FormValue("iEmail")
		m.BirthDay = r.FormValue("iBirthday")
		m.Notes = r.FormValue("iNotes")
		m.Status, _ = strconv.Atoi(r.FormValue("iStatus"))
		id := UpdMember(m)
		if mID == 0 {
			PrintUIHeader(w, "Saving Member")
			fmt.Fprintf(w, "<p>New Member saved, Name=%s --  Go to <a href='/members'>Search</a>.</p>",
				m.Name)
		} else {
			log.Printf("Data saved for the Member, ID=%d, Name=%s", id, m.Name)
			http.Redirect(w, r, "/members", http.StatusFound)
		}
	} else {
		NotAllowed(w, r)
	}
}
