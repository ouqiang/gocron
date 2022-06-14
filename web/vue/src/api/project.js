import httpClient from '../utils/httpClient'

export default {
    list(query, callback) {
        httpClient.get('/project', query, callback)
    },
    store(data, callback) {
        httpClient.post('/project', data, callback)
    },
    all(callback) {
        httpClient.get('/project/all', {}, callback)
    },
}