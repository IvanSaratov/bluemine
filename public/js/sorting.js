function customStat(status) {
    switch (status) {
        case "Новая": {
            return 0;
            break;
        }
        case "В процессе": {
            return 1;
            break;
        }
        case "Отложена": {
            return 2;
            break;
        }
        case "Закрыта": {
            return 3;
            break;
        }
    }
}

$(document).ready(function() {
    $('#tasks_exec #sort_stat_down').click()
    $('#tasks_creator #sort_stat_down').click()
})

var tableID;

$('#tasks_exec #sort_name_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        return $(a).find('#name a').text() < $(b).find('#name a').text()
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_exec  #sort_name_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        return $(a).find('#name a').text() > $(b).find('#name a').text()
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_creator #sort_name_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        return $(a).find('#name a').text() < $(b).find('#name a').text()
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_creator  #sort_name_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        return $(a).find('#name a').text() > $(b).find('#name a').text()
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_exec #sort_stat_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customStat($(a).find('#stat').text())
        bVal = customStat($(b).find('#stat').text())

        return aVal > bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_exec #sort_stat_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customStat($(a).find('#stat').text())
        bVal = customStat($(b).find('#stat').text())

        return aVal < bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_creator #sort_stat_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customStat($(a).find('#stat').text())
        bVal = customStat($(b).find('#stat').text())

        return aVal > bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_creator #sort_stat_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customStat($(a).find('#stat').text())
        bVal = customStat($(b).find('#stat').text())

        return aVal < bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_exec #sort_rate_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = parseInt($(a).find('#rate').text())
        bVal = parseInt($(b).find('#rate').text())

        return aVal > bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_exec #sort_rate_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = parseInt($(a).find('#rate').text())
        bVal = parseInt($(b).find('#rate').text())

        return aVal < bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_creator #sort_rate_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = parseInt($(a).find('#rate').text())
        bVal = parseInt($(b).find('#rate').text())

        return aVal > bVal
    }).appendTo('#' + tableID + ' tbody')
})

$('#tasks_creator #sort_rate_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = parseInt($(a).find('#rate').text())
        bVal = parseInt($(b).find('#rate').text())

        return aVal < bVal
    }).appendTo('#' + tableID + ' tbody')
})