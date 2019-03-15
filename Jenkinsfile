pipeline {
    agent {
        label 'docker'
    }
    stages {
        stage('Build') {
            when {
                not { buildingTag() }
            }
            steps {
                sh "docker build --build-arg VERSION=${BUILD_TAG} ."
            }
        }
        stage('Deploy') {
            when {
                tag pattern: "v\\d+\\.\\d+\\.\\d+(-\\w+-\\d+)?", comparator: "REGEXP"
            }
            steps {
                script {
                    VERSION = TAG_NAME[1..-1]
                }
                sh "echo Version is ${VERSION}"
                sh "docker build --tag ${GIT_COMMIT} --build-arg VERSION=${VERSION} ."
                sh "docker tag ${GIT_COMMIT} fint/graphql-cli:${VERSION}"
                sh "docker tag ${GIT_COMMIT} fint/graphql-cli:latest"
                withDockerRegistry([credentialsId: 'asgeir-docker', url: '']) {
                    sh "docker push fint/graphql-cli:${VERSION}"
                    sh "docker push fint/graphql-cli:latest"
                }
            }
        }
    }
}
