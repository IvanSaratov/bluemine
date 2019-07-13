$('#new_task').hide(0);
$('#new_tmpl').hide(0);
$('#new_group').hide(0);

$('#add_new_task').click(function() {
    $('#new_task').show(300);
    $('#task_send').show(300);
    $('#task_change').hide(0);
});

$('#add_new_tmpl').click(function() {
    $('#new_tmpl').show(300);
});

$('#add_new_group').click(function() {
    $('#new_group').show(300);
    $('#group_send').show(300);
    $('#group_change').hide(0);
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

$('#input_task_tmpl').on('change', function() {
    var ID = $(this).children(":selected").attr("id");

    if (ID != undefined) {
        $.get("/gettmpldata", {tmpl_id: ID}).done(function(data) {
            $('#input_task_stat').val(data.TmplStat)
            $('#input_task_priority').val(data.TmplPriority)
            $('#input_task_rate').val(data.TmplRate)
        })
    } else {
        $('#input_task_stat').val("Новая")
        $('#input_task_priority').val("Низкий")
        $('#input_task_rate').val(0)
    }
});

function makeAdmin(){
    var ID = $('.pagetitle').attr("id");

    $.ajax({
        url: "/makeadmin",
        method: "POST",
        data: {
            user_id: ID
        },
        success: function(){
            location.reload();
        }
    });
}
function removeAdmin(){
    var ID = $('.pagetitle').attr("id");

    $.ajax({
        url: "/removeadmin",
        method: "POST",
        data: {
            user_id: ID
        },
        success: function(){
            location.reload();
        }
    });
}

function taskAdd(){
    var name = document.getElementById("input_task_name").value;
    var desc = document.getElementById("input_desc").value;
    var stat = document.getElementById("input_task_stat").value;
    var priority = document.getElementById("input_task_priority").value;
    var exec = $('#input_task_exec :selected').attr('id');
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
                task_exec: exec,
                task_exec_type: exec_type,
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

function taskChange(){
    var name = document.getElementById("input_task_name").value;
    var desc = document.getElementById("input_desc").value;
    var stat = document.getElementById("input_task_stat").value;
    var priority = document.getElementById("input_task_priority").value;
    var exec = $('#input_task_exec :selected').attr('id');
    var exec_type = $('#input_task_exec :selected').attr('class');
    var rate = document.getElementById("input_task_rate").value;
    var date_start = document.getElementById("input_task_date_start").value;
    var date_end = document.getElementById("input_task_date_end").value;

    if (name.length == 0 || exec.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/tasks/change",
            method: "POST",
            data: {
                task_name: name,
                task_desc: desc,
                task_stat: stat,
                task_priority: priority,
                task_exec: exec,
                task_exec_type: exec_type,
                task_rate: rate,
                task_start: date_start,
                task_end: date_end
            },
            success: function(){
                location.reload();
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

function groupAdd() {
    var name = document.getElementById('input_group_name').value;
    var list = $('.user:checked').serialize();

    if (name.length == 0 || list.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/group/new",
            method: "POST",
            data: {
                group_name: name,
                user_list: list
            },
            success: function(){
                location.href = "/groups";
            }
        });
    }
};

function groupChange() {
    var href = location.href.split('/');
    var groupID = href[5];
    var name = document.getElementById('input_group_name').value;
    var list = $('.user:checked').serialize();

    if (name.length == 0 || list.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/group/change",
            method: "POST",
            data: {
                group_id: groupID,
                group_name: name,
                user_list: list
            },
            success: function(){
                location.reload();
            }
        });
    }
};