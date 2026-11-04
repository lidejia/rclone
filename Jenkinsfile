// The image is built and pushed by the buildImage step in the k3s-homelib
// shared library (vars/buildImage.groovy), which also tags it with the
// commit's short SHA.
pipeline {
    agent any

    stages {
        stage('Build & Push') {
            steps {
                buildImage(image: 'nexus3.workstation.lidejia.info/docker/rclone', dockerfile: 'Dockerfile.jenkins', tags: ['1.75.4-hotfix'])
            }
        }
    }
}
