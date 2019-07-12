function tChange() {
    var href = location.href.split('/');
    var id = href[5];
    if (id.length == 0) {
        alert("Пустое значение")
    } else {
        $.get("/gettaskdata", {task_id: id}).done(function(data) {
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