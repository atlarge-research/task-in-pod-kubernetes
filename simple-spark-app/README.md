# Simple spark application #
1. ```sbt compile```
2. ```sbt package```

From the same folder, run ```$SPARK_HOME/bin/spark-submit --class "SimpleApp" --master local target/scala-2.12/simpleapp_2.12-1.0.jar```
