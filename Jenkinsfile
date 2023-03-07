pipeline{
    agent any
    stages {
        stage('Build') {
            steps {
                echo '${BUILD_NUMBER}'
            }
        }
        stage('Test') {
            steps {
                echo '${BUILD_TAG}'
            }
        }
        stage('Deploy') {
            steps {
                echo '${JENKINS_URL}'
            }
        }
    }
}