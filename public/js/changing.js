function tChange() {
    var href = location.href.split('/');
    var id = href[5];
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
            $('#' + data.TaskExecutorID).prop('selected', true);
            $('#input_task_rate').val(data.TaskRate)
            $('#task_date_start').val(data.TaskDateStart)
            $('#task_date_end').val(data.TaskDateEnd)
        });
        $.get("/gettaskdesc", { id: id }).done(function(data){
            $('#input_desc').val(data)
        })
    }
}

function gChange() {
    var href = location.href.split('/');
    var id = href[5];
    $.get("/groupchange", { id: id }).done(function(data){
        $('#task_change').show(0);
        $('#task_send').hide(0);
        $('#new_group').show(300);
        $('#input_group_name').val(data.GroupName)

        for (i = 0; i < data.GroupMembers.length; i++) {
            $("input[value='"+data.GroupMembers[i].UserName+"']").prop('checked', true);
        };
    })
};
