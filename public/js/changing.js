function taskFillForChange() {
    var id = $('.pagetitle').attr("id");

    if (id.length == 0) {
        alert("Пустое значение")
    } else {
        $.get("/get/taskdata", {task_id: id}).done(function(data) {
            $('#task_change').show(0);
            $('#task_send').hide(0);
            $('#new_task').show(300);

            $('#input_task_name').val(data.TaskName)
            $('#input_task_stat').val(data.TaskStat)
            $('#input_task_priority').val(data.TaskPriority)
            $('#' + data.TaskExecutorName).prop('selected', true);
            $('#input_task_rate').val(data.TaskRate)
            $('#input_task_date_start').val(data.TaskDateStart)
            $('#input_task_date_end').val(data.TaskDateEnd)

            for (i = 0; i < data.TaskChecklist.length; i++) {
                var html = '<div id="check_container" style="width=100%;"><input class="checkbox" type="checkbox" name="checkbox" value="' + data.TaskChecklist[i].CheckName + '"><label for="' + data.TaskChecklist[i].CheckName + '">' + data.TaskChecklist[i].CheckName + '</label><span id="icoremovecheckbox" onclick="removeCheckbox(this)"></span></div>';
                $('#task_checklist_wrap').append(html);
                if (data.TaskChecklist[i].Checked) {
                    $("input[value='"+data.TaskChecklist[i].CheckName+"']").prop('checked', true);
                }
            };
        });
        var type = 'tasks';
        $.get("/get/taskdesc", { id: id, type: type }).done(function(data){
            $('#input_task_desc').val(data)
            MDParse(data, $('#markdown_output'))
        })
    }
}

function groupFillForChange() {
    var id = $('.pagetitle').attr("id");

    $.get("/change/group", { id: id }).done(function(data){
        $('#group_change').show(0);
        $('#group_send').hide(0);
        $('#new_group').show(300);
        $('#input_group_name').val(data.GroupName)

        for (i = 0; i < data.GroupMembers.length; i++) {
            $("input[value='"+data.GroupMembers[i].UserName+"']").prop('checked', true);
        };
    })
};

function taskChange(){
    var id = $('.pagetitle').attr("id");
    var name = document.getElementById("input_task_name").value;
    var desc = document.getElementById("input_task_desc").value;
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
            url: "/change/task",
            method: "POST",
            data: {
                task_id: id,
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
                location.reload();
            }
        });
    }
};

function groupChange() {
    var groupID = $('.pagetitle').attr("id");
    var name = document.getElementById('input_group_name').value;
    var list = $('.user:checked').serialize();

    if (name.length == 0 || list.length == 0) {
        alert("Пустое значение")
    } else {
        $.ajax({
            url: "/change/group",
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