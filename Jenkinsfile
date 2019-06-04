pipeline {
    agent {
        node {
            label 'mesos'
        }
    }
    stages {
        stage('test') {
            steps {
                sh 'make test'
            }
        }
    }
}
