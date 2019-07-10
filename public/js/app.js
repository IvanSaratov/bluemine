$('#new_task').hide(0);
$('#new_tmpl').hide(0);
$('#new_group').hide(0);

function setDate() {
    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1;
    var yyyy = today.getFullYear();

    if (dd < 10) {
        dd = '0'+dd
    } 

    if (mm < 10) {
        mm = '0'+mm
    } 

    today = yyyy + '-' + mm + '-' + dd;

    document.getElementById("input_task_date_start").value = today;
}

window.onload = function() {
    setDate();
};

$('#add_new_task').click(function() {
    $('#new_task').show(300);
});

$('#add_new_tmpl').click(function() {
    $('#new_tmpl').show(300);
});

$('#add_new_group').click(function() {
    $('#new_group').show(300);
});

$('#new_task .new_item_close').click(function() {
    $('#new_task').hide(300);
});

$('#new_tmpl .new_item_close').click(function() {
    $('#new_tmpl').hide(300);
});

$('#new_group .new_item_close').click(function() {
    $('#new_group').hide(300);
});

function taskAdd(){
    var name = document.getElementById("input_task_name").value;
    var desc = document.getElementById("input_task_desc").value;
    var stat = document.getElementById("input_task_stat").value;
    var priority = document.getElementById("input_task_priority").value;
    var exec = document.getElementById("input_task_exec").value;
    var exec_type = $('#input_task_exec :selected').attr('class');
    var rate = document.getElementById("input_task_rate").value;
    var date_start = document.getElementById("input_task_date_start").value;
    var date_end = document.getElementById("input_task_date_end").value;

    if (name.length == 0 || exec.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/tasks/new",
            method: "POST",
            data: {
                task_name: name,
                task_desc: desc,
                task_stat: stat,
                task_priority: priority,
                exec_name: exec,
                exec_type: exec_type,
                task_rate: rate,
                task_start: date_start,
                task_end: date_end
            },
            success: function(){
                location.href = "/tasks";
            }
        });
    }
};

function tmplAdd(){
    var name = document.getElementById("input_tmpl_name").value;
    var stat = document.getElementById("input_tmpl_stat").value;
    var priority = document.getElementById("input_tmpl_priority").value;
    var rate = document.getElementById("input_tmpl_rate").value;

    if (name.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/tmpl/new",
            method: "POST",
            data: {
                tmpl_name: name,
                tmpl_stat: stat,
                tmpl_priority: priority,
                tmpl_rate: rate,
            },
            success: function(){
                location.href = "/tasks";
            }
        });
    }
};

function gAdd() {
    var groupName = document.getElementById('input_group').value;
    var checkboxes = document.getElementsByName('users');
    var users = [];

    for (var i=0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            users.push(checkboxes[i].value);
        }
    }

    if (groupName.length == 0 || checkboxes.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/group/new",
            method: "POST",
            data: {
                input_group: groupName,
                user_list: users.toString()
            },
            success: function(){
                location.href = "/groups";
            }
        });
    }
};