function tAdd(){
    var name = document.getElementById("inputtaskname").value;
    var desc = document.getElementById("inputdesc").value;
    var exec = document.getElementById("inputexec").value;
    var exec_type = $('#inputexec :selected').attr('class');
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
                exec_type: exec_type
            },
            success: function(){
                location.href = "/tasks";
            }
        });
    }
};