pipeline {
    agent any

    environment {
        IMAGE = 'nexus3.workstation.lidejia.info/docker/rclone'
    }

    stages {
        stage('Build & Push') {
            steps {
                script {
                    def tag = env.GIT_COMMIT.take(7)
                    container('kaniko') {
                        sh "/kaniko/executor --context=dir://${WORKSPACE} --dockerfile=Dockerfile.jenkins --destination=${IMAGE}:1.75.4-hotfix --destination=${IMAGE}:${tag} --ignore-path=/product_uuid"
                    }
                }
            }
        }
    }
}
