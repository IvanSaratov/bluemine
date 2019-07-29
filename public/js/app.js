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

$(document).keyup(function(e) {
    if (e.key === "Escape") {
        var items = $('.new_item');
        for (var i = 0; i < items.length; i++) {
            var itemID = $(items[i]).attr('id');
            var item = $('#' + itemID)
            if (item.css('display') == 'block') {
                item.hide(300);
            }
        }
    }
});

$(document).ready(function() {
    $('#iconew').hover(
        function() {
            $('#new_item').css('display', 'block');
        },
        function() {
            $('#new_item').css('display', 'none');
        }
    )
});

window.onload = function() {
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
};

$('#input_task_tmpl').on('change', function() {
    var ID = $(this).children(":selected").attr("id");

    if (ID != undefined) {
        $.get("/get/tmpldata", {tmpl_id: ID}).done(function(data) {
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

$(document).ready(function() {
    $('#input_desc').bind('input propertychange', function() {
        var markdown = $('#input_desc').val();
        MDParse(markdown, $('#markdown_output'))
    })
})

$(document).ready(function() {
    if ($('#task_checklist .checkbox').length == 0) {
        $('#task_checklist').hide(0);
    }
})

function makeAdmin(){
    var ID = $('.pagetitle').attr("id");

    $.ajax({
        url: "/admin/make",
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
        url: "/admin/remove",
        method: "POST",
        data: {
            user_id: ID
        },
        success: function(){
            location.reload();
        }
    });
}