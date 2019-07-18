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

function customPriority(priority) {
    switch (priority) {
        case "Низкий": {
            return 0;
            break;
        }
        case "Средний": {
            return 1;
            break;
        }
        case "Высокий": {
            return 2;
            break;
        }
    }
}

$(document).ready(function() {
    $('.sort_stat_down').click()
})

var tableID;

$('.sort_name_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#name a').text().toUpperCase()
        bVal = $(b).find('#name a').text().toUpperCase()

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_name_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#name a').text().toUpperCase()
        bVal = $(b).find('#name a').text().toUpperCase()

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_stat_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customStat($(a).find('#stat').text())
        bVal = customStat($(b).find('#stat').text())

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_stat_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customStat($(a).find('#stat').text())
        bVal = customStat($(b).find('#stat').text())

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_pr_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customPriority($(a).find('#pr').text())
        bVal = customPriority($(b).find('#pr').text())

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_pr_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = customPriority($(a).find('#pr').text())
        bVal = customPriority($(b).find('#pr').text())

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_exec_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#exec a').text().toUpperCase()
        bVal = $(b).find('#exec a').text().toUpperCase()

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_exec_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#exec a').text().toUpperCase()
        bVal = $(b).find('#exec a').text().toUpperCase()

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_add_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#date_start').text()
        bVal = $(b).find('#date_start').text()

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_add_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#date_start').text()
        bVal = $(b).find('#date_start').text()

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_limit_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#date_end').text()
        bVal = $(b).find('#date_end').text()

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_limit_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#date_end').text()
        bVal = $(b).find('#date_end').text()

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_rate_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = parseInt($(a).find('#rate').text())
        bVal = parseInt($(b).find('#rate').text())

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_rate_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = parseInt($(a).find('#rate').text())
        bVal = parseInt($(b).find('#rate').text())

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_update_down').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#last_update').text()
        bVal = $(b).find('#last_update').text()

        return (aVal < bVal) ? -1 : (aVal > bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})

$('.sort_update_up').click(function() {
    tableID = $(this).closest('table').attr('id');
    $('#' + tableID + ' .sort_icon').css('color', 'white')
    $(this).css('color', 'rgb(54, 161, 248)')

    $('#' + tableID + ' tbody tr').sort(function (a, b) {
        aVal = $(a).find('#last_update').text()
        bVal = $(b).find('#last_update').text()

        return (aVal > bVal) ? -1 : (aVal < bVal) ? 1 : 0;
    }).appendTo('#' + tableID + ' tbody')
})