function tChange() {
    var id = $('.pagetitle').attr("id");

    if (id.length == 0) {
        alert("Пустое значение")
    } else {
        $.get("/gettaskdata", {task_id: id}).done(function(data) {
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
        });
        var type = 'tasks';
        $.get("/gettaskdesc", { id: id, type: type }).done(function(data){
            $('#input_desc').val(data)
        })
    }
}

function gChange() {
    var id = $('.pagetitle').attr("id");

    $.get("/group/change", { id: id }).done(function(data){
        $('#group_change').show(0);
        $('#group_send').hide(0);
        $('#new_group').show(300);
        $('#input_group_name').val(data.GroupName)

        for (i = 0; i < data.GroupMembers.length; i++) {
            $("input[value='"+data.GroupMembers[i].UserName+"']").prop('checked', true);
        };
    })
};
