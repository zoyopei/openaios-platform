# [cartpole] 快速构建强化学习智能体完成倒立摆小游戏
 
## 你会获得什么？
本次任务提供了一个构建强化学习智能体的案例，你将会在5分钟的时间内，快速学习并完成强化学习教材经典标配案例“倒立摆”小游戏的智能体构建。

## 背景
倒立摆（CartPole V0）是强化学习教材中的必修题。倒立摆体验地址：https://fluxml.ai/experiments/cartPole/ 

游戏里面有一个小车，它的上面竖着一根杆子，小车可以无摩擦的左右运动，每次重置后的初始状态会有所不同。小车需要左右移动来保持杆子竖直


![cartpole](https://img-blog.csdnimg.cn/20190426152813188.gif)


## 操作指南

### 运行帮助

- JupyterLab使用
  - 打开代码文件之后，我们看到任务栏中有如下图标，选中文件中的一个模块，点击三角形按钮就可以执行该模块，如果要提前停止，点击正方形按钮。
  - 如果想要重新执行代码，点击圆形箭头，然后从一个模块开始执行

![2](https://ftp.bmp.ovh/imgs/2021/05/be7c96c5d42ffcd4.png)

### 运行demo
一共分为4部分，分别是初始化程序，配置参数，训练和预测。按顺序执行每个cell即可完成比赛。

- 第一部分为初始化过程，会安装强化学习相关的依赖，期待看到的结果
![](https://ftp.bmp.ovh/imgs/2021/05/90bc47d91b313db8.png)
*图1 初始化*

- 第二部分为配置参数，可以根据代码注释的理解适当调整参数，也可以直接执行，期待看到的输出如图

![](https://ftp.bmp.ovh/imgs/2021/05/4d943df17374626d.png)
*图2 配置参数*


- 第三部分为训练过程，执行时间略长，期待的输出如图
![](https://ftp.bmp.ovh/imgs/2021/05/08a896548403eda5.png)
*图3 训练*

- 第四部分为预测部分，可以看到模拟的过程，以及最终的结果，期待的输出如图
![](https://ftp.bmp.ovh/imgs/2021/05/6c625db78de847a4.png) *图4 预测*

### 评分方法

执行完看到图4视为完成，请将图4截图，通过邮件提交至opensource@4paradigm.com，邮件名称为“闯关赛结果截图提交” 

## 进一步了解游戏

游戏的每个状态由4个变量确定
小车位置 : Cart Position: [-4.8, 4.8]
小车速度 : Cart Velocity: [-Inf, Inf]
杆子的角度 : Pole Angle: [-24 degree, 24 degree]
杆子的速度 : Pole Velocity at Tip: [-Inf, Inf]

每次有两个操作：
0: 表示给小车施加向左的力
1: 表示给小车施加向右的力

注：施加的力大小是固定的，但减小或增大的速度不是固定的，它取决于当时杆子与竖直方向的角度。角度不同，产生的速度和位移也不同。


## 设计模拟器以及奖励
### 奖励机制:
每一步都给出1的奖励，包括终止状态。

### 初始状态:
初始状态所有变量都从[-0.05,0.05]中随机取值。

### 达到下列条件之一片段结束
杆子与竖直方向角度超过12度
小车位置距离中心超过2.4（小车中心超出画面）
片段长度超过200
