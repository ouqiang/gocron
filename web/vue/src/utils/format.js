export default {
    formatDatetime(date) {
        let padStart = function (i) {
            return i.toString().padStart(2, '0')
        }
        let d = new Date(date)
        if (d.getFullYear() === 1) {
            return ''
        }
        return d.getFullYear() + '-' + padStart(d.getMonth() + 1) + '-' + padStart(d.getDate()) + ' '
            + padStart(d.getHours()) + ':' + padStart(d.getMinutes()) + ':' + padStart(d.getSeconds())
    }
}