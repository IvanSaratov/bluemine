function tAdd(){
    var name = document.getElementById("input_task_name").value;
    var desc = document.getElementById("input_desc").value;
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