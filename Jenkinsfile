env.name                   = "github-backup"
env.description            = "github-backup is a small tool to backup all private and public repositories from a specific GitHub organization."
env.maintainer             = "ops <ops@productsup.com>"
env.homepage               = "https://github.com/productsupcom/github-backup"
String dockerImage         = "golang:1.21"
env.version
env.branch
env.gitCommitHash
env.gitCommitAuthor
env.gitCommitMessage
env.package_file_name

pipeline {
    agent { label 'jenkins-4'}
    options {
        buildDiscarder(
            logRotator(
                numToKeepStr: '5',
                artifactNumToKeepStr: '5'
            )
        )
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        disableConcurrentBuilds()
        skipDefaultCheckout()
    }

    stages {
        stage("Checkout") {
            steps {
                gitCheckout()
            }
        }

        stage('Prepare Info') {
            steps {
                prepareInfo()
            }
        }

        // Pull docker image to use for the tests / build
        stage('Pull image') {
            steps {
                script {
                    docker.withRegistry('https://docker.productsup.com', 'docker-registry') {
                        sh "docker pull ${dockerImage}"
                    }
                }
            }
        }

        stage('Run Tests') {
            agent {
                docker {
                    image "${dockerImage}"
                    // reuseNode is needed to make sure we got all the required information
                    reuseNode true
                }
            }
            steps {
                sh 'go test'
            }
        }

        // Build go binary
        stage('Build go package') {
            agent {
                docker {
                    image "${dockerImage}"
                    // reuseNode is needed to make sure we got all the required information
                    reuseNode true
                }
            }
            steps {
                // build and set internal variable appVersion to current version
                sh "go build -o ./build/${name} -ldflags \"-X main.version=${env.version}\""
            }
        }

        // build production deb package when we build a tag.
        // we use different naming for the packages in dev and prod
        stage ('Build and publish prod deb package') {
            when {
                buildingTag()
            }
            steps {
                setPackageName(customName: "${env.name}_${env.version}")
                // build
                buildDebPackageBin(
                    package_internal_name: "${env.name}",
                    package_file_name: "${env.package_file_name}",
                    version: "${env.version}",
                    description: "${env.description}",
                    homepage: "${env.homepage}",
                    maintainer: "${env.maintainer}"
                )
                // publish
                publishDebPackage(package_name: "${env.package_file_name}_all.deb")
                uploadFileToGithubRelease('productsupcom', env.name, env.TAG_NAME, "${env.name}", "build/${env.name}")
            }
        }
    }

    // Run post jobs steps
    post {
        cleanup {
            cleanWs deleteDirs: true
        }
    }
}
