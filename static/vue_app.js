const App = {
    data() {
        return {
            title: 'Тесты',

            testList: undefined,
            test: undefined,

            access_token: undefined,
            loginOpen: false,
            loginShowPassword: false,
            loginInfo: {
                username: "demo_account",
                password: "qwerty1234",
                fingerprint: undefined,
            },
            registrationOpen: false,
            registrationShowPassword: false,
            registrationInfo: {
                username: "demo_account",
                password: "qwerty1234",
                name: undefined,
                email: "example@example.com",
            },
        }
    },
    // https://v3.ru.vuejs.org/ru/guide/instance.html#диаграмма-жизненного-цикла
    beforeCreate() {},
    created() {
        // Initialize the agent at application startup.
        const fpPromise = import('https://openfpcdn.io/fingerprintjs/v3')
        .then(FingerprintJS => FingerprintJS.load())

        // Get the visitor identifier when you need it.
        fpPromise
        .then(fp => fp.get())
        .then(result => {
            this.loginInfo.fingerprint = result.visitorId // result.components
            this.refreshToken()
        })
        .catch(error => console.error(error))

        if (location.pathname === '/') {
            this.retrieveAllTests()
        }

        addEventListener('popstate', (event) => { 
            if (location.pathname === '/') {
                this.openTestsList()
            }
        })
    },
    beforeMount() {},
    mounted() {},
    beforeUnmount() {},
    unmounted() {},

    beforeUpdate(){}, // управление путями сделать?
    updated(){},

    computed: {},
    watch: {},
    methods: {
        page(){
            return location.pathname == '/' ? 'list' : 'test'
        },

        parseJWT (token) {
            const base64Payload = token.split('.')[1]
            const base64 = base64Payload.replace(/-/g, '+').replace(/_/g, '/')
            const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(c => {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
            }).join(''))

            return JSON.parse(jsonPayload)
        },

        refreshToken() {
            // проверить локально что не закончился JWT токен и если просрочен, то обновить
            const URL = '/api/auth/refresh-token'
            fetch(URL, {
                method: 'POST',
                headers: {
                    'Host': '',
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    'fingerprint': this.loginInfo.fingerprint
                })
            }).then(response => {
                if (response.status === 200)
                    response.json().then(json => {
                        this.access_token = `Bearer ${json.access_token}`
                        this.title = "Тесты"
                        history.pushState(null, null, '/')

                        const parsedJWT = this.parseJWT(json.access_token)
                        const expiredAt = new Date(parsedJWT.exp)
                        const now = (new Date()).getTime() / 1000
                        const timeout = (expiredAt - now - 3) * 1000
                        setTimeout(() => {
                            if (this.access_token === `Bearer ${json.access_token}`) {
                                this.refreshToken()
                                /*if (this.access_token === `Bearer ${json.access_token}`) {
                                    this.access_token = null
                                }*/
                            }
                        }, timeout)
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        openTestsList() {
            this.title = 'Тесты'
            history.pushState(null, null, '/')
            this.retrieveAllTests()
        },

        signUpOpenForm() {
            this.registrationShowPassword = false
            this.registrationOpen = true
            this.$nextTick(() => {
                this.$refs.reg.focus()
            })
        },

        signInOpenForm() {
            this.loginShowPassword = false
            this.loginOpen = true
            this.$nextTick(() => {
                this.$refs.login.focus()
            })
        },

        signIn() {
            if (this.loginInfo.username == '' || this.loginInfo.password == '') {
                console.log("sign-in: required login and password")
                return
            }

            const URL = '/api/auth/sign-in'
            fetch(URL, {
                method: 'POST',
                //mode: 'cors',
                credentials: 'same-origin', // include, *same-origin, omit
                headers: {
                    'Host': '',
                    'Content-Type': 'application/json',
                },
                redirect: 'error',
                body: JSON.stringify(this.loginInfo)
            }).then(response => {
                if (response.status === 200)
                    response.json().then(json => {
                        this.access_token = `Bearer ${json.access_token}`
                        this.loginOpen = false

                        const parsedJWT = this.parseJWT(json.access_token)
                        const expiredAt = new Date(parsedJWT.exp)
                        const now = (new Date()).getTime() / 1000
                        const timeout = (expiredAt - now) * 1000
                        // обнулить просроченный токен
                        setTimeout(() => {
                            if (this.access_token === `Bearer ${json.access_token}`) {
                                this.access_token = null
                                //this.refreshToken()
                            }
                        }, timeout)
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        signUp() {
            if (this.registrationInfo.username == '' || this.registrationInfo.password == '' || this.registrationInfo.email == '') {
                console.log("error")
                return
            }

            const URL = '/api/auth/sign-up'
            fetch(URL, {
                method: 'POST',
                headers: {
                    'Host': '',
                    'Content-Type': 'application/json',
                },
                redirect: 'error',
                body: JSON.stringify(this.registrationInfo)
            }).then(response => {
                if (200 <= response.status && response.status < 300)
                    response.json().then(json => {
                        this.registrationOpen = false
                        this.loginOpen = true
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        signOut() {
            const URL = '/api/auth/sign-out'
            fetch(URL, {
                method: 'POST',
                headers: {
                    'Host': '',
                },
            }).then(response => {
                if (response.status === 200) {
                    this.retrieveAllTests()
                    this.access_token = undefined
                    this.title = 'Тесты'
                    history.pushState(null, null, '/')
                } else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        retrieveAllTests() {
            const URL = '/api/test'
            fetch(URL, {
                method: 'GET',
                headers: {
                    'Host': '',
                }
            }).then(response => {
                if (response.status === 200)
                    response.json().then(json => {
                        this.testList = json
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        retrieveTestById(test_id) {
            if (!this.access_token) {
                this.signInOpenForm()
                return
            }

            const URL = `/api/test/${test_id}/full`
            fetch(URL, {
                method: 'GET',
                headers: {
                    'Host': '',
                    'credentials': 'omit',
                    'Authorization': this.access_token,
                }
            }).then(response => {
                if (response.status === 200)
                    response.json().then(json => {
                        console.log("retrieveTestById", json)
                        this.test = json
                        this.testList = undefined
                        this.title = this.test.test.title
                        history.pushState(null, null, 'test/' + test_id)

                        for (let i = 0; i < this.test.questions.length; i++) {
                            let question = this.test.questions[i]
                            if (question.answer_type === 'manySelect') {
                                question.response = []
                            }

                            for (let j = 0; j < this.test.answers.length; j++) {
                                let answer = this.test.answers[j]
                                if (question.id === answer.question_id && 0 < answer.answer.length) {
                                    question.response = question.answer_type === 'manySelect' ? answer.answer : answer.answer[0]
                                }
                            }

                            // true_answers подсветить зелёным
                            for (let j = 0; j < question.show_answers.length; j++) {
                                question.show_answers[j] = {
                                    value: question.show_answers[j],
                                    style: question.true_answers && question.true_answers.includes(question.show_answers[j]) ? { color: "green" } : {}
                                }
                            }

                            //if (question.answer_type === 'freeField') {}
                        }
                        console.log(this.test)
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        setAnswer(event) {
            if (!this.access_token) {
                signInOpenForm()
                return
            }

            const question_index = Number(event.target.getAttribute('question_index'))
            const question = this.test.questions[question_index]

            if (!question.response || question.response == '' || question.response.length === 0) {
                console.log('empty answer')
                return
            }

            const URL = '/api/answer'
            fetch(URL, {
                method: 'POST',
                headers: {
                    'Host': '',
                    'Authorization': this.access_token,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "question_id": question.id,
                    "answer": question.answer_type === 'manySelect' ? question.response : [question.response],
                })
            }).then(response => {
                if (response.status === 200)
                    response.json().then(json => {
                        console.log(json)
                        //this.test.questions.push(question)
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },

        completeTest(test_id) {
            if (!this.access_token) {
                this.signInOpenForm()
                return
            }

            const URL = `/api/test-answer`
            fetch(URL, {
                method: 'POST',
                headers: {
                    'Host': '',
                    'Authorization': this.access_token,
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "test_id": test_id
                })
            }).then(response => {
                if (response.status === 200)
                    response.json().then(json => {
                        this.retrieveTestById(test_id)
                    })
                else
                    console.log(response)
            }).catch(err => console.log(err))
        },
    },
}

Vue.createApp(App).mount('#vue-app')
