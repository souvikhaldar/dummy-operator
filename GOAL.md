**a8s Applicant Homework: Write and Deploy a Custom Kubernetes Controller**

Thegoalofthishomeworkistousethe [OperatorSDK ](https://github.com/operator-framework/operator-sdk)towriteasmallKubernetesCustomControllerin GoandthendeployitonaKubernetescluster.

Tocompletethehomeworkyoucanfreelyuseanyresourceyouwant. ThatincludesGoogle,StackOver- flow,anymailinglist/chat(e.g.Slack)thatyouthinkcouldhelpyou,etc…. However,youMUSTdothe homeworkonyourownwithoutanyassitancefromfriends,relativesetc….

Aspartofthehomework,you’llhavetowritecodeandotherartifacts. Everythingyouproduceshouldbe committedinthisrepo.

**Step1(Optional): FamiliarizeyourselfwiththeConceptsandTools**

**Thisstephasnodeliverableandishereonlytohelpyou. Feelfreetoskipit.**

Tostartwith,youshoulgetanunderstandingofwhataKubernetescontrollerisandthebasicsofwriting onewiththeOperatorSDKandGo. Followsomeresources,boththeoreticalandpractical,thatshould helpyou. Obviously,youcanalsodoyourownresearchandconsultresourcesthatarenotlistedhereif thathelpsyou.

- [KubernetesAPIObjects](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/)
- [KubernetesControllers](https://kubernetes.io/docs/concepts/architecture/controller/)
- [KubernetesOperators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [OperatorSDKinstallation](https://sdk.operatorframework.io/docs/installation/)
- [OperatorSDKtutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)

**Step2: DefineaCustomResourceandaControllerthatLogstheSpecoftheCustom Resource**

UsetheOperatorSDKto:

1. CreateanewcustomAPItypecalled“Dummy”.
1. Createanewcustomcontrollerforresourcesofkind“Dummy”.

“Dummy”resourcesshouldhaveno status fieldanda spec fieldthathasonlyonestringsubfieldcalled “message”. Followsanexample:

apiVersion**:** interview.com/v1alpha1 kind **:** Dummy

metadata**:**

name**:** dummy1

namespace**:** default

spec**:**

message**:** "I'm just a dummy"

ThecustomcontrollermustprocesseachDummyAPIobjectsimplybyloggingitsname,namespaceand thevalueof spec.message.

**Step3: GiveDummiesaStatusandMaketheCustomControllerWritetoit**

ExtendtheDummycustomAPItypebygivingtoita status fieldwhichcontainsonlyastringsubfield called specEcho.

ThecustomcontrollershouldnotonlylogeachDummyname,namespaceand spec.message,itshould nowalsocopythevalueof spec.message into status.specEcho.

Forexample,afteraDummyhasbeenprocessedbythecustomcontrolleritshouldlooklike:

apiVersion**:** interview.com/v1alpha1 kind **:** Dummy

metadata**:**

name**:** dummy1

namespace**:** default

spec**:**

message**:** "I'm just a dummy" status**:**

specEcho**:** "I'm just a dummy"

**Step4: AssociateaPodtoeachDummyAPIobject**

ExtendthecustomcontrollertoalsomakeitcreateaPodforeachDummyAPIobjectthatiscreated. WhenaDummyisdeleteditsPodshouldalsoceasetoexist(eitherimmediatelyorshortlyafterwards). ThePodshouldrun [nginx.](https://hub.docker.com/_/nginx)

Futhermore,theDummycustomAPItypeshouldbeextendedbygivingtoitastatusfieldthattracksthe statusofthePodassociatedtotheDummy.

Forexample,aDummyshouldinitiallystartas:

apiVersion**:** interview.com/v1alpha1 kind **:** Dummy

metadata**:**

name**:** dummy1

namespace**:** default

spec**:**

message**:** "I'm just a dummy" status**:**

specEcho**:** "I'm just a dummy" podStatus**:** "Pending"

ButonceitsPodisupandrunningitshouldbecome:

apiVersion**:** interview.com/v1alpha1 kind **:** Dummy

metadata**:**

name**:** dummy1

namespace**:** default

spec**:**

message**:** "I'm just a dummy" status**:**

specEcho**:** "I'm just a dummy" podStatus**:** "Running"

**Step5: RuntheCustomControlleronKubernetes**

Deploy the custom controller on a Kubernetes cluster (for example, one created with [Minikube) ](https://minikube.sigs.k8s.io/docs/start/)and testit. Alsowriteamarkdownfileforthereviewerofyourhomeworkwithinstructionsonhowtodo thesame. Noticethatyouwillneedtopushthecontainerimageofyourcustomcontrollertoapublicly accessiblecontainerregistry(forexampleon [DockerHub),](https://hub.docker.com/)sothatthereviewercantestyoursolution.

We organized the homework in incremental steps to help you do one thing at a time, but you’re not constrainedtofollowthestepswedescribed: youcandothehomeworkinanyorder/wayyoulike. In fact,weonlycareaboutthefinalresultratherthantheresultaftereachintermediatestep.

Goodluck!
3
d
