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
    var checklist = '';

    $('#new_task .checkbox').each(function() {
        checklist += $(this).val() + '=' + $(this).prop('checked') + '&';
    });
    checklist = checklist.slice(0, -1);

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
                task_end: date_end,
                task_checklist: checklist
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

function wikiAdd() {
    var name = document.getElementById('input_wiki_name').value;
    var article = document.getElementById('input_desc').value;
    var father_id = $('#input_wiki_father :selected').attr('id');

    if (name.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/wiki/new",
            method: "POST",
            data: {
                wiki_name: name,
                article: article,
                father_id: father_id
            },
            success: function(){
                location.href = "/wiki";
            }
        });
    }
};