;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is part of the IWs Project: https://bitbucket.org/rsucasas/iws
;;
;;
;;
;; Copyright: Roi Sucasas Font 2017.
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.data_db_test
  (:require [clojure.test :refer :all]
            [atos.class.logs.logs :as LOG]
            [atos.class.db.db :refer :all]))


(deftest test-01
  (LOG/info "> Testing iws.clj.core.db.db [test-01] ...")
  (testing "Testing iws.clj.core.db.db [test-01] ... "
    (let [mdb2 (create-mem-db)]
      (dbw mdb2 [:db-A :table-A1 :content] {:id 1 :txt "0100"})
      (dbw mdb2 [:db-A :table-A2 :content] {:id 2 :txt "0200"})
      (is (= (dbr mdb2 [:db-A :table-A2 :content :txt]) "0200")))))


(deftest test-02
  (LOG/info "> Testing iws.clj.core.db.db [test-02] ...")
  (testing "Testing iws.clj.core.db.db [test-02] ... "
    (let [mdb2 (open-db "/tmp/storef")]
      (dbw mdb2 [:db-A :table-A1 :content] {:id 1 :txt "0100"})
      (dbw mdb2 [:db-A :table-A2 :content] {:id 2 :txt "0234"})
      (save-mem-table-to-disk mdb2 "/tmp/store" :db-A)
      (load-mem-table-from-disk  mdb2 "/tmp/store" :db-A :db-B)
      (remove-db "/tmp/store")
      (is (= (dbr mdb2 [:db-B :table-A2 :content :txt]) "0234")))))


(deftest test-03
  (LOG/info "> Testing iws.clj.core.db.db [test-03] ...")
  (testing "Testing iws.clj.core.db.db [test-03] ... "
    (let [mdb2 (create-mem-db)]
      (doseq [i (range 20)]
        (dbw mdb2 [:db-A (keyword (str "table-A" i "")) :content] {:id 1 :txt "0100"}))
      (is (= (count (dbr mdb2 [:db-A])) 20)))))


(deftest test-04
  (LOG/info "> Testing iws.clj.core.db.db [test-04] ...")
  (testing "Testing iws.clj.core.db.db [test-04] ... "
    (let [mdb2 (open-db "/tmp/storef")]
      (is (= (dbr mdb2 [:db-A :table-A1 :content]) {:id 1 :txt "0100"}))
      (is (= (dbr mdb2 [:db-A :table-A2 :content]) {:id 2 :txt "0234"}))
      (dbd mdb2 [:db-A :table-A2 :content])
      (is (= (dbr mdb2 [:db-A :table-A2 :content]) {})))))


(deftest test-05
  (LOG/info "> Testing iws.clj.core.db.db [test-05] ...")
  (testing "Testing iws.clj.core.db.db [test-05] ... "
    (let [mdb2 (open-db "/tmp/storef")]
      (dbw mdb2 [:tab1 :reg1 :id1] {:name "name" :ports [1 2]})
      (is (= (dbr mdb2 [:tab1 :reg1 :id1]) {:name "name" :ports [1 2]}))

      (dbu mdb2 [:tab1 :reg1 :id1] {:name "name2" :ports [1 2]})
      (is (= (dbr mdb2 [:tab1 :reg1 :id1]) {:name "name2" :ports [1 2]}))
      (is (= (dbr mdb2 [:tab1 :reg1 :id1 :name]) "name2"))

      (dbu mdb2 [:tab1 :reg1 :id1 :name] "nameXXXX")
      (is (= (dbr mdb2 [:tab1 :reg1 :id1 :name]) "nameXXXX")))))


(deftest test-06
  (LOG/info "> Testing iws.clj.core.db.db [test-06] ...")
  (testing "Testing iws.clj.core.db.db [test-06] ... "
    (let [mdb2 (open-db "/tmp/storef")]
      (LOG/info ">> (dbr mdb2 [:tab1 :reg1 :id1]) = " (dbr mdb2 [:tab1 :reg1 :id1]))
      (LOG/info ">> (dbr mdb2 [:tab1 :reg1]) = " (dbr mdb2 [:tab1 :reg1]))
      (LOG/info ">> (dbr mdb2 [:tab1]) = " (dbr mdb2 [:tab1])))

    (let [mdb2 (open-db "/tmp/storef")]
      (dbw mdb2 [:tab2 :reg2] {:id2 {:name "name" :ports [1 2]}})
      (dbw mdb2 [:tab2 :reg3] {:id4 {:name "name4" :ports [1 2 3 4]}})
      (LOG/info ">> (dbr mdb2 [:tab2 :reg2]) = " (dbr mdb2 [:tab2 :reg2]))
      (LOG/info ">> (dbr mdb2 [:tab2 :reg3]) = " (dbr mdb2 [:tab2 :reg3]))

      (dba mdb2 [:tab2 :reg2] {:id3 {:name "name2" :ports [1 2 3]}})
      (LOG/info ">> (dbr mdb2 [:tab2 :reg2]) = " (dbr mdb2 [:tab2 :reg2]))

      (LOG/info ">> (dbr mdb2 [:tab2]) = " (dbr mdb2 [:tab2]))

      (delete-table-db mdb2 :tab2)

      (LOG/info ">> (dbr mdb2 [:tab2]) = " (dbr mdb2 [:tab2]) " nil? :)")
      (is (= (dbr mdb2 [:tab2]) nil)))

    (remove-db "/tmp/storef")))


(deftest test-07
  (LOG/info "> Testing iws.clj.core.db.db [test-07] ...")
  (testing "Testing iws.clj.core.db.db [test-07] ... "
    (let [mdb2 (open-db "/tmp/storef")]
      (dbw mdb2 [:tasks :id-task-1]
                { :name         "dname"
                  :deployment   "dname"
                  :service      ["sname11", "sname12"]
                  :volumenClaim ["vname11", "vname12"]
                  :ports        [1 2 3]
                  :status       "deleted" })
      (LOG/info ">> [:tasks :id-task-1] =     " (dbr mdb2 [:tasks :id-task-1]))

      (LOG/info ">> updating [:tasks :id-task-1] =     " (dbu mdb2 [:tasks :id-task-1 :status] "unknown"))

      (LOG/info ">> [:tasks :id-task-1] =     " (dbr mdb2 [:tasks :id-task-1]))

      (dbw mdb2 [:tasks :id-task-2]
                { :name         "dname2"
                  :deployment   "dname2"
                  :service      ["sname21", "sname22"]
                  :volumenClaim ["vname21", "vname22"]
                  :ports        [1 2 3 4]
                  :status       "running" })

      (LOG/info ">> all tasks =     " (dbr mdb2 [:tasks]))
      (LOG/info ">> running tasks = " (into {} (filter (fn [[k v]] (= (v :status) "running")) (dbr mdb2 [:tasks]))))

      (dbd mdb2 [:tasks :id-task-1])
      (LOG/info ">> all tasks =     " (dbr mdb2 [:tasks]))

      (dbw mdb2 [:tasks :id-task-1]
                { :name         "dname"
                  :deployment   "dname"
                  :service      ["sname11", "sname12"]
                  :volumenClaim ["vname11", "vname12"]
                  :ports        [1 2 3]
                  :status       "deleted" })
      (LOG/info ">> all tasks =     " (dbr mdb2 [:tasks])))
    (remove-db "/tmp/storef")))



(defn test-ns-hook []
  (test-01)
  (test-02)
  (test-03)
  (test-04)
  (test-05)
  (test-06)
  (test-07))
