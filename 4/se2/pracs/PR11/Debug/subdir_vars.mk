################################################################################
# Automatically-generated file. Do not edit!
################################################################################

# Add inputs and outputs from these tool invocations to the build variables 
CFG_SRCS += \
../app.cfg 

CMD_SRCS += \
../TM4C123GH6PM.cmd 

LIB_SRCS += \
../driverlib.lib 

C_SRCS += \
../computos.c \
../main.c 

GEN_CMDS += \
./configPkg/linker.cmd 

GEN_FILES += \
./configPkg/linker.cmd \
./configPkg/compiler.opt 

GEN_MISC_DIRS += \
./configPkg 

C_DEPS += \
./computos.d \
./main.d 

GEN_OPTS += \
./configPkg/compiler.opt 

OBJS += \
./computos.obj \
./main.obj 

GEN_MISC_DIRS__QUOTED += \
"configPkg" 

OBJS__QUOTED += \
"computos.obj" \
"main.obj" 

C_DEPS__QUOTED += \
"computos.d" \
"main.d" 

GEN_FILES__QUOTED += \
"configPkg/linker.cmd" \
"configPkg/compiler.opt" 

C_SRCS__QUOTED += \
"../computos.c" \
"../main.c" 


