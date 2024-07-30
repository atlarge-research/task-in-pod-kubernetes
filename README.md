# Task-in-Pod Scheduling Support for Kubernetes and Apache Spark Stack
This is a bachelor's project of the Vrije Universiteit Amsterdam student. The project is composed of two parts:
<ol>
  <li>Task-aware scheculer</li>
  <li>Comparison of frameworks for Operator creation</li>
</ol>

## 1. Task-aware scheduler ##
A task-aware scheduler is a prototype of an external scheduler added to the Spark framework. It is made as a Kubernetes Operator and modifies the original workflow of the Spark-Kubernetes application submission process.

The *Task-aware scheduler* folder contains a modified version of Apache Spark source code and the code for the actual scheduler (folder *spark-operator*). A README file there describes how to create a distribution and run a Spark application with novel functionality.

## 2. Frameworks comparison ##
To choose which frameworks to use, we have conducted an experiment with the three most commonly used frameworks: Kopf, Metacontroller, and Operator SDK. The workflow and all the required files are presented in the *frameworks* folder.
