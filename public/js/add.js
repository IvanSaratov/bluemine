function tAdd(){
    var name = document.getElementById("input_task_name").value;
    var desc = document.getElementById("input_desc").value;
    var stat = document.getElementById("input_stat").value;
    var priority = document.getElementById("input_priority").value;
    var exec = document.getElementById("input_exec").value;
    var exec_type = $('#input_exec :selected').attr('class');
    var rate = document.getElementById("input_rate").value;
    var date_start = document.getElementById("input_date_start").value;
    var date_end = document.getElementById("input_date_end").value;
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