$('#new_task').hide(0);
$('#new_tmpl').hide(0);
$('#new_group').hide(0);

$('#add_new_task').click(function() {
    $('#new_task').show(300);
    $('#task_send').show(0);
    $('#task_change').hide(0);
});

$('#add_new_tmpl').click(function() {
    $('#new_tmpl').show(300);
});

$('#add_new_group').click(function() {
    $('#new_group').show(300);
    $('#group_send').show(0);
    $('#group_change').hide(0);
});

$(document).ready(function() {
    $('legend').click(function() {
        if ($('#filters').hasClass('show')) {
            $({deg: 0}).animate({deg: -90}, {
                duration: 200,
                step: function(now){
                    $('#icoarrowright').css({
                        transform: "rotate(" + now + "deg)"
                    });
                }
            })
            $('#filters').hide(200)
            $('#filters').removeClass('show')
        } else {
            $({deg: 0}).animate({deg: 360}, {
                duration: 200,
                step: function(now){
                    $('#icoarrowright').css({
                        transform: "rotate(" + now + "deg)"
                    });
                }
            })
            $('#filters').show(200)
            $('#filters').addClass('show')
        }
    })
    $('legend').click()
})

function showIt(prop, filter) {
    $('#tasks tbody tr').each(function() {
        if ($(this).find('#' + prop).text() == filter) {
            $(this).show(200);
        }
    })
}

function hideIt(prop, filter) {
    $('#tasks tbody tr').each(function() {
        if ($(this).find('#' + prop).text() == filter) {
            $(this).hide(200);
        }
    })
}

$(document).ready(function() {
    $('#filters').find('.record :checkbox').change(function() {
        if (this.checked) {
            showIt($(this).parent('.record').attr('id'), $('label[for="' + $(this).attr('name') + '"]').text())
        } else {
            hideIt($(this).parent('.record').attr('id'), $('label[for="' + $(this).attr('name') + '"]').text())
        }
    })
})

$(document).ready(function() {
    $('#filters').find('.record :checkbox').each(function() {
        if (!this.checked) {
            hideIt($(this).parent('.record').attr('id'), $('label[for="' + $(this).attr('name') + '"]').text())
        }
    })
})

$(document).keyup(function(e) {
    switch (e.key) {
        case "1": {
            if ($('#new_tmpl').css('display') != 'block' && $('#new_group').css('display') != 'block') {
                $('#new_task').show(300);
                $('#task_send').show(0);
                $('#task_change').hide(0);
            }
            break;
        }
        case "2": {
            if ($('#new_task').css('display') != 'block' && $('#new_group').css('display') != 'block') {
                $('#new_tmpl').show(300);
            }
            break;
        }
        case "3": {
            if ($('#new_task').css('display') != 'block' && $('#new_tmpl').css('display') != 'block') {
                $('#new_group').show(300);
                $('#group_send').show(0);
                $('#group_change').hide(0);
            }
            break;
        }
    }
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

$(document).on('click', '#icowikihide', function() {
    var wikiid = $(this).parent('.item').closest('.wiki_art').attr('id')
    if (wikiid == undefined) {
        wikiid = $(this).parent('.nested_item').closest('.wiki_art').attr('id')
    }
    $('#child_of_' + wikiid).hide(200)
    $(this).attr('id', 'icowikishow')
})
$(document).on('click', '#icowikishow', function() {
    var wikiid = $(this).parent('.item').closest('.wiki_art').attr('id')
    if (wikiid == undefined) {
        wikiid = $(this).parent('.nested_item').closest('.wiki_art').attr('id')
    }
    $('#child_of_' + wikiid).show(200)
    $(this).attr('id', 'icowikihide')
})

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
            $('#' + data.TmplExecType).prop('selected', true);
            $('#input_task_rate').val(data.TmplRate)
        })
    } else {
        $('#input_task_stat').val("В процессе")
        $('#input_task_priority').val("Низкий")
        $('#input_task_rate').val(0)
    }
});

$(document).ready(function() {
    $('#input_task_desc').bind('input propertychange', function() {
        var markdown = $('#input_task_desc').val();
        MDParse(markdown, $('#markdown_output'))
    })
})

$(document).ready(function() {
    $('.group_user').sort(function (a, b) {
        aVal = $(a).find('label').text().toUpperCase()
        bVal = $(b).find('label').text().toUpperCase()
    
        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo($('.group_user').parent('.record'))
})

$(document).ready(function() {
    if ($('#checklist .checkbox').length == 0) {
        $('#checklist').hide(0);
    } else {
        $('#date_end').parent('#task_info_left').append('<div class="record" id="progress"><div class="label">Готовность:</div><div class="value"><div id="bar"><div id="bar_status"></div></div></div><div id="percent_number">0%</div></div>')
        var width = 0;
        var bar = $('#bar_status');
        var checkedPercent = $('#checklist .checkbox:checked').length / $('#checklist .checkbox').length * 100;
        var interval = setInterval(frame, 10);

        function frame() {
            if (width >= checkedPercent) {
                clearInterval(interval);
            } else {
                width++;
                bar.css('width', width + '%');
                $('#percent_number').text(width + '%');
            }
        }
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